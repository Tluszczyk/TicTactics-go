package database

import (
	"errors"
	"testing"
)

func TestParseDatabaseDeploymentOption(t *testing.T) {
	tests := []struct {
		name             string
		deploymentOption string
		want             DatabaseDeploymentOption
		wantErr          error
	}{
		{
			name:             "Valid deployment option",
			deploymentOption: "DYNAMO",
			want:             DYNAMO,
			wantErr:          nil,
		},
		{
			name:             "Invalid deployment option",
			deploymentOption: "INVALID",
			want:             NONE,
			wantErr:          errors.ErrUnsupported,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDatabaseDeploymentOption(tt.deploymentOption)
			if err != tt.wantErr {
				t.Errorf("ParseDatabaseDeploymentOption() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseDatabaseDeploymentOption() = %v, want %v", got, tt.want)
			}
		})
	}
}
