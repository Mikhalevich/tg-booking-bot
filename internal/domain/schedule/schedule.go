package schedule

import (
	"context"
	"fmt"
	"strings"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
)

type schedule struct {
	repository port.ScheduleRepository
	sender     port.MessageSender
}

func New(
	repository port.ScheduleRepository,
	sender port.MessageSender,
) *schedule {
	return &schedule{
		repository: repository,
		sender:     sender,
	}
}

func (s *schedule) GetAllTemplates(ctx context.Context, info port.MessageInfo) error {
	tmpls, err := s.repository.GetAllTemplates(ctx)
	if err != nil {
		return fmt.Errorf("get all templates: %w", err)
	}

	if len(tmpls) == 0 {
		if err := s.sender.ReplyText(
			ctx,
			info.ChatID,
			info.MessageID,
			"no schedule templates",
		); err != nil {
			return fmt.Errorf("reply to message empty: %w", err)
		}

		return nil
	}

	for _, tmpl := range tmpls {
		if err := s.sender.ReplyTextMarkdown(
			ctx,
			info.ChatID,
			info.MessageID,
			s.convertScheduleToString(tmpl),
		); err != nil {
			return fmt.Errorf("reply message template: %w", err)
		}
	}

	return nil
}

func (s *schedule) convertScheduleToString(tmpl port.ScheduleTemplate) string {
	var (
		sortedHoursByDay = []string{
			fmt.Sprintf("*Name:* %s", s.sender.EscapeMarkdown(tmpl.Name)),
			fmt.Sprintf("*Description:* %s", s.sender.EscapeMarkdown(tmpl.Description)),
			"*Schedule:*",
		}
		weekDays = [...]port.DayOfWeek{port.Mon, port.Tue, port.Wed, port.Thu, port.Fri, port.Sat, port.Sun}
	)

	for _, weekDay := range weekDays {
		hours := convertHoursForDayOfWeekToString(tmpl.Schedule.WorkingHours, weekDay)
		if hours != "" {
			sortedHoursByDay = append(sortedHoursByDay, fmt.Sprintf("%s: %s", weekDay.String(), hours))
		}
	}

	return strings.Join(sortedHoursByDay, "\n")
}

func convertHoursForDayOfWeekToString(workingHours []port.HoursByDay, targetDay port.DayOfWeek) string {
	for _, wh := range workingHours {
		for _, day := range wh.Days {
			if day != targetDay {
				continue
			}

			var output string
			for _, hours := range wh.Hours {
				output += fmt.Sprintf("*%s\\-%s* ", hours.Start.String(), hours.End.String())
			}

			return output
		}
	}

	return ""
}
