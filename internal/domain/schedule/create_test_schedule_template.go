package schedule

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/msginfo"
)

func (s *schedule) CreateTestScheduleTemplate(ctx context.Context, info msginfo.Info) error {
	if err := s.repository.CreateTestScheduleTemplate(ctx, port.ScheduleTemplate{
		Name:        "test_schedule_template_name_1",
		Description: "test_schedule_templaet_description_1",
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
	); err != nil {
		return fmt.Errorf("create test schedule template: %w", err)
	}

	return nil
}
