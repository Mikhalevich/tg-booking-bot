package schedule

import (
	"context"
	"fmt"
	"strings"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
)

func (s *schedule) GetAllTemplates(ctx context.Context, info port.MessageInfo) error {
	tmpls, err := s.repository.GetAllTemplates(ctx)
	if err != nil {
		return fmt.Errorf("get all templates: %w", err)
	}

	if len(tmpls) == 0 {
		s.sender.ReplyText(ctx, info.ChatID, info.MessageID, "no schedule templates")
		return nil
	}

	for _, tmpl := range tmpls {
		s.sender.ReplyTextMarkdown(ctx, info.ChatID, info.MessageID, s.convertScheduleToString(tmpl))
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
