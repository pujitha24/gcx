package definitions_test

import (
	"encoding/json"
	"testing"

	"github.com/grafana/gcx/internal/providers/slo/definitions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuery_PreservesSourceDatasourceUID(t *testing.T) {
	tests := []struct {
		name string
		in   string
		get  func(definitions.Query) string
	}{
		{
			name: "freeform",
			in:   `{"type":"freeform","freeform":{"query":"up","sourceDatasourceUid":"source-ds"}}`,
			get:  func(q definitions.Query) string { return q.Freeform.SourceDatasourceUID },
		},
		{
			name: "ratio",
			in:   `{"type":"ratio","ratio":{"successMetric":{"prometheusMetric":"a","type":"counter"},"totalMetric":{"prometheusMetric":"b","type":"counter"},"sourceDatasourceUid":"source-ds"}}`,
			get:  func(q definitions.Query) string { return q.Ratio.SourceDatasourceUID },
		},
		{
			name: "threshold",
			in:   `{"type":"threshold","threshold":{"thresholdExpression":"x","threshold":{"value":1,"operator":">"},"sourceDatasourceUid":"source-ds"}}`,
			get:  func(q definitions.Query) string { return q.Threshold.SourceDatasourceUID },
		},
		{
			name: "failureThreshold",
			in:   `{"type":"failureThreshold","failureThreshold":{"failureThresholdExpression":"x","threshold":{"value":1,"operator":">"},"sourceDatasourceUid":"source-ds"}}`,
			get:  func(q definitions.Query) string { return q.FailureThreshold.SourceDatasourceUID },
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var q definitions.Query
			require.NoError(t, json.Unmarshal([]byte(tt.in), &q))
			assert.Equal(t, "source-ds", tt.get(q))

			out, err := json.Marshal(q)
			require.NoError(t, err)
			assert.Contains(t, string(out), `"sourceDatasourceUid":"source-ds"`)
		})
	}
}
