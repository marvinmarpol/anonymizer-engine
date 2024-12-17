package service

import (
	"context"
	"strings"

	"github.com/marvinmarpol/golang-boilerplate/internal/pkg/utils/cryptho"
	"github.com/marvinmarpol/golang-boilerplate/internal/protection/query"
	"github.com/sirupsen/logrus"
)

func (s *ServiceImpl) traverseMapAndDecrypt(ctx context.Context, data interface{}, keyPrefixes []string, valuePrefixes []string, decryptPrefix string) {
	switch vv := data.(type) {
	case []interface{}:
		for i, item := range vv {
			switch itemTyped := item.(type) {
			case map[string]interface{}:
				// Recurse into nested maps
				s.traverseMapAndDecrypt(ctx, itemTyped, keyPrefixes, valuePrefixes, decryptPrefix)

			case []interface{}:
				// Recurse into nested slices
				s.traverseMapAndDecrypt(ctx, itemTyped, keyPrefixes, valuePrefixes, decryptPrefix)

			default:
				if strVal, ok := item.(string); ok && strings.HasPrefix(strVal, decryptPrefix) {
					// decrypt to original value
					original, err := s.decryptMask(ctx, strVal)
					if err != nil {
						continue
					}

					// assign encrypted value
					vv[i] = original
				}

			}
		}

	case map[string]interface{}:
		for key, value := range vv {
			switch v := value.(type) {
			case map[string]interface{}:
				s.traverseMapAndDecrypt(ctx, v, keyPrefixes, valuePrefixes, decryptPrefix)

			case []interface{}:
				// Traverse slices
				for i, item := range v {
					if nestedMap, ok := item.(map[string]interface{}); ok {
						s.traverseMapAndDecrypt(ctx, nestedMap, keyPrefixes, valuePrefixes, decryptPrefix)
					} else {
						// Check if the value in the slice matches any of the value prefixes
						if strVal, ok := item.(string); ok && strings.HasPrefix(strVal, decryptPrefix) {
							// decrypt to original value
							original, err := s.decryptMask(ctx, strVal)
							if err != nil {
								continue
							}

							v[i] = original
						}
					}
				}

			default:
				if strVal, ok := value.(string); ok && strings.HasPrefix(strVal, decryptPrefix) {
					// decrypt to original value
					original, err := s.decryptMask(ctx, strVal)
					if err != nil {
						continue
					}

					// assign encrypted value
					vv[key] = original
				}

			}
		}
	}

}

func (s *ServiceImpl) decryptMask(ctx context.Context, strVal string) (string, error) {
	token := ""
	splits := strings.SplitN(strVal, decryptPrefix, 2)
	if len(splits) == 2 {
		token = splits[1]
	}

	// get row from mask table by token
	maskI, err := s.Queries.GetMaskQuery.Handle(ctx, query.GetMaskQuery{Token: token})
	if err != nil {
		logrus.WithContext(ctx).WithField("err", err).Error("Failed to get mask")
		return "", err
	}

	// decrypt the encrypted key
	maskT := maskI.Result
	originalKey, _ := cryptho.RsaDecrypt(s.PriKey, maskT.Key)

	// decrypt to original value
	original, err := cryptho.AESDecrypt(maskT.Cypher, originalKey)
	if err != nil {
		logrus.WithContext(ctx).WithField("err", err).Error("Failed to decrypt")
	}

	return original, err
}
