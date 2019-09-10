package repository

import (
	"fmt"
	shared "shared-bob"
)

type BobRepository struct {
	requestHelper *shared.RequestHelper
}

func NewBobRepository(bobUrl string) *BobRepository {
	return &BobRepository{
		requestHelper: shared.NewRequestHelper(bobUrl),
	}
}

func (br *BobRepository) Sync() error {
	resp, err := br.requestHelper.Post("/api/v1/sync", nil, nil)
	if err != nil {
		return err
	}

	if resp != nil && resp.StatusCode != 200 {
		return fmt.Errorf("bob responded with code %d", resp.StatusCode)
	}

	return nil
}