package service

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/marvinmarpol/golang-boilerplate/internal/pkg/utils/cryptho"
	"github.com/marvinmarpol/golang-boilerplate/internal/protection/command"
	"github.com/marvinmarpol/golang-boilerplate/internal/protection/domain/mask"
	"github.com/marvinmarpol/golang-boilerplate/internal/protection/entity"
	"github.com/marvinmarpol/golang-boilerplate/internal/protection/query"
	"github.com/sirupsen/logrus"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type ServiceImpl struct {
	Commands command.Commands
	Queries  query.Queries
	PubKey   *rsa.PublicKey
	PriKey   *rsa.PrivateKey
}

func NewServiceImpl(cmd command.Commands, q query.Queries, pubKey *rsa.PublicKey, priKey *rsa.PrivateKey) *ServiceImpl {
	return &ServiceImpl{cmd, q, pubKey, priKey}
}

func (s *ServiceImpl) Deidentify(ctx context.Context, cmd interface{}) (interface{}, error) {
	s.traverseMapAndEncrypt(ctx, cmd, keyPrefixes, valuePrefixes, decryptPrefix)
	return cmd, nil
}
func (s *ServiceImpl) Reidentify(ctx context.Context, cmd interface{}) (interface{}, error) {
	s.traverseMapAndDecrypt(ctx, cmd, keyPrefixes, valuePrefixes, decryptPrefix)
	return cmd, nil
}
func (s *ServiceImpl) GetCypher(ctx context.Context, cmd entity.GetCypherPayload) (interface{}, error) {
	return s.Queries.GetCypherQuery.Handle(ctx, query.GetCypherQuery{Token: cmd.Token})
}

func (s *ServiceImpl) RotateKeys(ctx context.Context, cmd entity.RotatePayload) (interface{}, error) {
	// init required vars
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	limit := cmd.BatchSize
	offset := 0
	counter := 0
	async := 0

	// loop until reach max or no more result
	for {
		candidates, err := s.Queries.GetRotateCandidateQuery.Handle(ctx, query.GetRotateCandidateQuery{
			Limit:         limit,
			Offset:        offset,
			DayDifference: cmd.DayDifference,
		})
		if err != nil {
			return counter, err
		}

		// rotate key of the retrieved masks
		for _, item := range candidates.Result {
			wg.Add(1)

			// Protect async counter increment with a mutex
			mu.Lock()
			async++
			mu.Unlock()

			go s.rotateKey(context.Background(), item, &wg, &async, &mu)

			// Check if we reached the max async limit
			mu.Lock()
			if async >= cmd.MaxAsyncProcess {
				mu.Unlock() // Unlock before waiting
				wg.Wait()   // Wait for goroutines to finish
				async = 0   // Reset async counter after waiting
			} else {
				mu.Unlock()
			}
		}
		time.Sleep(time.Duration(cmd.MsDelayEachJob) * time.Millisecond)

		// count retrieved row, break if reach max or no data retrieved
		counter += len(candidates.Result)
		if len(candidates.Result) < 1 || counter > cmd.Max {
			break
		}

		// set offset for next rows
		offset += candidates.Limit
	}

	return counter, nil
}

func (s *ServiceImpl) rotateKey(ctx context.Context, item mask.Mask, wg *sync.WaitGroup, asyncCounter *int, mu *sync.Mutex) error {
	defer wg.Done()

	originalKey, err := cryptho.RsaDecrypt(s.PriKey, item.Key)
	if err != nil {
		logrus.WithContext(ctx).WithField("err", err).Error("Failed to get original key")
		return err
	}
	plainText, err := cryptho.AESDecrypt(item.Cypher, originalKey)
	if err != nil {
		logrus.WithContext(ctx).WithField("err", err).Error("Failed to get plaintext")
		return err
	}
	newKey, err := gonanoid.Generate(possibleChars, encKeyLength)
	if err != nil {
		logrus.WithContext(ctx).WithField("err", err).Error("Failed to create new key")
		return err
	}
	newCypher, err := cryptho.AESEncrypt(plainText, newKey)
	if err != nil {
		logrus.WithContext(ctx).WithField("err", err).Error("Failed to create new cypher")
		return err
	}
	newEncKey, err := cryptho.RsaEncrypt(s.PubKey, newKey)
	if err != nil {
		logrus.WithContext(ctx).WithField("err", err).Error("Failed to create new encryption key")
		return err
	}

	// update mask with new key and cypher
	s.Commands.UpdateMaskCommand.Handle(context.Background(), command.UpdateMaskCommand{
		Token:  item.Token,
		Key:    newEncKey,
		Cypher: newCypher,
	})

	// Protect asyncCounter using mutex
	mu.Lock()
	*asyncCounter--
	mu.Unlock()

	return nil
}

// retry function
func RetryUntilSuccess[T any](fn func() (T, error), maxAttempts int) (T, error) {
	var result T
	var err error

	for i := 0; i < maxAttempts; i++ {
		result, err = fn()
		if err == nil {
			return result, nil // Success, no error
		}

		logrus.Error(fmt.Sprintf("Error: %v Counter: %v", err, i+1))
	}

	return result, errors.New("max attempts reached without success")
}
