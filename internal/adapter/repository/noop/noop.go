package noop

import (
	"context"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
)

var _ port.ScheduleRepository = (*noop)(nil)

type noop struct {
}

func New() *noop {
	return &noop{}
}

//nolint:funlen
func (n *noop) GetAllTemplates(ctx context.Context) ([]port.ScheduleTemplate, error) {
	return []port.ScheduleTemplate{
		{
			Name:        "test_schedule_name_1",
			Description: "test_schedule_description_1",
			Schedule: port.Schedule{
				WorkingHours: []port.HoursByDay{
					{
						Days: []port.DayOfWeek{port.Mon, port.Tue, port.Wed, port.Thu, port.Fri},
						Hours: []port.TimeInterval{
							{
								Start: port.Time{
									Hour: 8,
								},
								End: port.Time{
									Hour: 12,
								},
							},
							{
								Start: port.Time{
									Hour: 13,
								},
								End: port.Time{
									Hour: 17,
								},
							},
						},
					},
				},
			},
		},
		{
			Name:        "test_schedule_name_2",
			Description: "test_schedule_description_2",
			Schedule: port.Schedule{
				WorkingHours: []port.HoursByDay{
					{
						Days: []port.DayOfWeek{port.Mon, port.Wed},
						Hours: []port.TimeInterval{
							{
								Start: port.Time{
									Hour: 8,
								},
								End: port.Time{
									Hour: 12,
								},
							},
							{
								Start: port.Time{
									Hour: 13,
								},
								End: port.Time{
									Hour: 17,
								},
							},
						},
					},
					{
						Days: []port.DayOfWeek{port.Fri},
						Hours: []port.TimeInterval{
							{
								Start: port.Time{
									Hour: 8,
								},
								End: port.Time{
									Hour: 12,
								},
							},
							{
								Start: port.Time{
									Hour: 13,
								},
								End: port.Time{
									Hour: 16,
								},
							},
						},
					},
				},
			},
		},
	}, nil
}
