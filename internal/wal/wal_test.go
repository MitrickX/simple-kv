package wal

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/MitrickX/simple-kv/internal/config"
	utilsOs "github.com/MitrickX/simple-kv/internal/utils/os"
	utilsTime "github.com/MitrickX/simple-kv/internal/utils/time"
)

type deps struct {
	os   utilsOs.OS
	time utilsTime.Time
}

func TestWAL_Write(t *testing.T) {
	testTime, _ := time.Parse(time.RFC3339, "2025-09-21T12:19:39+03:00")

	tests := []struct {
		name       string
		config     config.ConfigWAL
		deps       func() deps
		queries    []string
		wantErrors []error
	}{
		{
			name:   "write_one_query",
			config: config.Default().WAL,
			deps: func() deps {
				return deps{
					os:   utilsOs.NewMockOS(t),
					time: utilsTime.NewMockTime(t),
				}
			},
			queries: []string{
				"PUT test 1",
			},
			wantErrors: []error{nil},
		},
		{
			name:   "write_two_queries",
			config: config.Default().WAL,
			deps: func() deps {
				return deps{
					os:   utilsOs.NewMockOS(t),
					time: utilsTime.NewMockTime(t),
				}
			},
			queries: []string{
				"PUT test 1",
				"PUT test2 2",
			},
			wantErrors: []error{nil, nil},
		},
		{
			name: "flush_open_file_error",
			config: config.ConfigWAL{
				FlushingBatchSize: 3,
				MaxSegmentSize:    config.DataSize(10 * config.MB),
				DataDirectory:     "test",
			},
			deps: func() deps {
				expectedFileName := testTime.Format(fileNameFromNowTimeLayout)

				mockTime := utilsTime.NewMockTime(t)
				mockTime.EXPECT().Now().Return(testTime)

				mockOs := utilsOs.NewMockOS(t)
				mockOs.EXPECT().OpenFile("test/"+expectedFileName, os.O_APPEND|os.O_TRUNC, os.FileMode(0644)).
					Return(nil, errors.New("open_file_error"))

				return deps{
					os:   mockOs,
					time: mockTime,
				}
			},
			queries: []string{
				"PUT test 1",
				"PUT test2 2",
				"PUT test3 3",
			},
			wantErrors: []error{
				nil,
				nil,
				fmt.Errorf("fail to open new wal segemnt file (%s): %w",
					testTime.Format(fileNameFromNowTimeLayout),
					errors.New("open_file_error"),
				),
			},
		},
		{
			name: "flush_write_to_file_error",
			config: config.ConfigWAL{
				FlushingBatchSize: 3,
				MaxSegmentSize:    config.DataSize(10 * config.MB),
				DataDirectory:     "test",
			},
			deps: func() deps {
				expectedFileName := testTime.Format(fileNameFromNowTimeLayout)

				mockTime := utilsTime.NewMockTime(t)
				mockTime.EXPECT().Now().Return(testTime)

				bufStr := `PUT test 1
PUT test2 2
PUT test3 3`

				file := utilsOs.NewMockFile(t)
				file.EXPECT().Write([]byte(bufStr)).Return(len(bufStr), nil)
				file.EXPECT().Sync().Return(errors.New("sync_error"))
				file.EXPECT().Name().Return(expectedFileName)

				mockOs := utilsOs.NewMockOS(t)
				mockOs.EXPECT().OpenFile("test/"+expectedFileName, os.O_APPEND|os.O_TRUNC, os.FileMode(0644)).
					Return(file, nil)

				return deps{
					os:   mockOs,
					time: mockTime,
				}
			},
			queries: []string{
				"PUT test 1",
				"PUT test2 2",
				"PUT test3 3",
			},
			wantErrors: []error{
				nil,
				nil,
				fmt.Errorf("fail to sync wal segment file (%s): %w",
					testTime.Format(fileNameFromNowTimeLayout),
					errors.New("sync_error"),
				),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dps := tt.deps()
			wal := NewWAL(tt.config, dps.os, dps.time)

			for i, q := range tt.queries {
				err := wal.Write(q)
				wantErr := tt.wantErrors[i]

				if wantErr != nil {
					if err == nil {
						t.Errorf("unexpected error: %v", err)
					} else if wantErr.Error() != err.Error() {
						t.Errorf("got error: %s, want: %s", err.Error(), wantErr.Error())
					}

					continue
				}

				if err != nil {
					t.Errorf("expect error: %v", wantErr)
				}
			}
		})
	}
}
