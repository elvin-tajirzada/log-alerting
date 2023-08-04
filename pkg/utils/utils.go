package utils

import "github.com/elvin-tacirzade/log-alerting/pkg/models"

func GroupBy(lokiLogEntry *models.LokiLogEntry) map[string]map[string][]models.Stream {
	groups := make(map[string]map[string][]models.Stream)

	for _, result := range lokiLogEntry.Data.Result {
		if groups["device"] == nil {
			groups["device"] = make(map[string][]models.Stream)
		}
		groups["device"][result.Stream.Dt] = append(groups["device"][result.Stream.Dt], result.Stream)

		if groups["path"] == nil {
			groups["path"] = make(map[string][]models.Stream)
		}
		groups["path"][result.Stream.Path] = append(groups["path"][result.Stream.Path], result.Stream)

		if groups["IP"] == nil {
			groups["IP"] = make(map[string][]models.Stream)
		}
		groups["IP"][result.Stream.IP] = append(groups["IP"][result.Stream.IP], result.Stream)
	}

	return groups
}
