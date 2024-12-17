package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/marvinmarpol/golang-boilerplate/internal/pkg/utils/cryptho"
	"github.com/marvinmarpol/golang-boilerplate/internal/pkg/utils/tuple"
	"github.com/marvinmarpol/golang-boilerplate/internal/protection/command"
	"github.com/marvinmarpol/golang-boilerplate/internal/protection/query"
	"github.com/sirupsen/logrus"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

func (s *ServiceImpl) traverseMapAndEncrypt(ctx context.Context, data interface{}, keyPrefixes []string, valuePrefixes []string, decryptPrefix string) {
	switch vv := data.(type) {
	case []interface{}:
		for i, item := range vv {
			switch itemTyped := item.(type) {
			case map[string]interface{}:
				// Recurse into nested maps
				s.traverseMapAndEncrypt(ctx, itemTyped, keyPrefixes, valuePrefixes, decryptPrefix)

			case []interface{}:
				// Recurse into nested slices
				s.traverseMapAndEncrypt(ctx, itemTyped, keyPrefixes, valuePrefixes, decryptPrefix)

			default:
				if strVal, ok := item.(string); ok && tuple.HasPrefixInList(&strVal, &valuePrefixes, true) && strings.TrimSpace(strVal) != "" {
					token, err := s.createMask(ctx, strVal)
					if err != nil {
						continue
					}

					// assign encrypted value
					vv[i] = decryptPrefix + token
				}
			}
		}

	case map[string]interface{}:
		for key, value := range vv {
			switch v := value.(type) {
			case map[string]interface{}:
				// Recurse into nested maps
				s.traverseMapAndEncrypt(ctx, v, keyPrefixes, valuePrefixes, decryptPrefix)

			case []interface{}:
				// Traverse slices
				for i, item := range v {
					if nestedMap, ok := item.(map[string]interface{}); ok {
						s.traverseMapAndEncrypt(ctx, nestedMap, keyPrefixes, valuePrefixes, decryptPrefix)
					} else {
						// Check if the value in the slice matches any of the value prefixes
						if strVal, ok := item.(string); ok && tuple.HasPrefixInList(&strVal, &valuePrefixes, true) && strings.TrimSpace(strVal) != "" {
							token, err := s.createMask(ctx, strVal)
							if err != nil {
								continue
							}

							// assign encrypted value
							v[i] = decryptPrefix + token
						}
					}
				}

			default:
				// Check if key has any of the specified prefixes and encrypt its value
				if tuple.HasPrefixInList(&key, &keyPrefixes, false) {
					strVal := fmt.Sprintf("%v", value) // Convert value to string if needed
					if strings.TrimSpace(strVal) == "" {
						continue
					}

					token, err := s.createMask(ctx, strVal)
					if err != nil {
						continue
					}

					// assign encrypted value
					vv[key] = decryptPrefix + token
				}

				// Check if value is a string and matches any of the value prefixes
				if strVal, ok := value.(string); ok && tuple.HasPrefixInList(&strVal, &valuePrefixes, true) && strings.TrimSpace(strVal) != "" {
					token, err := s.createMask(ctx, strVal)
					if err != nil {
						continue
					}

					// assign encrypted value
					vv[key] = decryptPrefix + token
				}
			}
		}
	}

}

func (s *ServiceImpl) createMask(ctx context.Context, strVal string) (string, error) {
	token, err := gonanoid.Generate(possibleChars, len(strVal))
	if err != nil {
		logrus.WithContext(ctx).WithField("err", err).Error("Failed to generate token")
		return "", err
	}
	hash, _ := cryptho.GenerateHash(strVal, cryptho.SHA256)
	encKey, _ := gonanoid.Generate(possibleChars, encKeyLength)
	encryptedValue, err := cryptho.AESEncrypt(strVal, encKey)
	encKeyKey, _ := cryptho.RsaEncrypt(s.PubKey, encKey)
	if err != nil {
		logrus.WithContext(ctx).WithField("err", err).Error("Failed to encrypt")
		return "", err
	}

	// attempt to create mask
	cmd := command.CreateMaskCommand{
		Token:  token,
		Hash:   hash,
		Key:    encKeyKey,
		Cypher: encryptedValue,
	}
	_, err = s.Commands.CreateMaskCommand.Handle(ctx, cmd)
	// update token until no unique constraint error
	if err != nil && strings.Contains(err.Error(), pgErrConstrationID) {
		if strings.Contains(err.Error(), hashConstraintKey) {
			oldToken, err := s.Queries.GetTokenQuery.Handle(ctx, query.GetTokenQuery{Hash: hash})
			if err != nil {
				logrus.WithContext(ctx).WithField("err", err).Error("Failed to get old token")
				return "", err
			}
			token = oldToken.Result
		} else {
			fn := func() (interface{}, error) {
				token, err = gonanoid.Generate(possibleChars, len(strVal))
				if err != nil {
					return nil, err
				}

				// retry create mask with new token
				cmd.Token = token
				return s.Commands.CreateMaskCommand.Handle(ctx, cmd)
			}

			_, err := RetryUntilSuccess(fn, maxRetry)
			if err != nil {
				logrus.WithContext(ctx).WithField("ErrorRetry:", err).Error("Failed to retryupdate")
				return "", err
			}
		}
	}

	return token, nil
}
