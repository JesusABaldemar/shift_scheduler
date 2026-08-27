package main

import (
	"fmt"
	"encoding/json"
)

type Schedule struct {
	Approved  bool

	Monday    [60]bool
	Tuesday   [60]bool
	Wednesday [60]bool
	Thursday  [60]bool
	Friday    [60]bool
}

func createSchedule() *Schedule {
	newSchedule := Schedule{
		Approved:  false,
		Monday:    [60]bool{},
		Tuesday:   [60]bool{},
		Wednesday: [60]bool{},
		Thursday:  [60]bool{},
		Friday:    [60]bool{},
	}
	return &newSchedule
}

// Function to select a time slot in the schedule, the time slot has to be at least 18 slots long, and the start and end slots have to be within the range of 0-48. If the time slot is invalid, the function will return without making any changes to the schedule.
// The time slot can't be longer than 52 slots either. If a invalid day is passed, the function will return without making any changes to the schedule.
func selectTimeSlot(schedule *Schedule, day string, startSlot int, endSlot int) {
	if startSlot < 0 || endSlot > 48 || startSlot >= endSlot || (endSlot-startSlot) < 18  || (endSlot-startSlot) > 52 {
		return
	}
	switch day {
	case "Monday":
		for i := startSlot; i < endSlot; i++ {
			schedule.Monday[i] = true
		}
	case "Tuesday":
		for i := startSlot; i < endSlot; i++ {
			schedule.Tuesday[i] = true
		}
	case "Wednesday":
		for i := startSlot; i < endSlot; i++ {
			schedule.Wednesday[i] = true
		}
	case "Thursday":
		for i := startSlot; i < endSlot; i++ {
			schedule.Thursday[i] = true
		}
	case "Friday":
		for i := startSlot; i < endSlot; i++ {
			schedule.Friday[i] = true
		}
	case "default":
		return
	}

}

// Function to clear a time slot in the schedule, the start and end slots have to be within the range of 0-60. If the time slot is invalid, the function will return without making any changes to the schedule.
// If the day is invalid, the function will return without making any changes to the schedule.
func clearTimeSlot(schedule *Schedule, day string, startSlot int, endSlot int) {
	if startSlot < 0 || endSlot > 60 || startSlot >= endSlot {
		return
	}
	switch day {
	case "Monday":
		for i := startSlot; i < endSlot; i++ {
			schedule.Monday[i] = false
		}
	case "Tuesday":
		for i := startSlot; i < endSlot; i++ {
			schedule.Tuesday[i] = false
		}
	case "Wednesday":
		for i := startSlot; i < endSlot; i++ {
			schedule.Wednesday[i] = false
		}
	case "Thursday":
		for i := startSlot; i < endSlot; i++ {
			schedule.Thursday[i] = false
		}
	case "Friday":
		for i := startSlot; i < endSlot; i++ {
			schedule.Friday[i] = false
		}
	case "default":
		return
	}

	
}

func clearSchedule(schedule *Schedule) {
		for i := 0; i < 60; i++ {
			schedule.Monday[i] = false
			schedule.Tuesday[i] = false
			schedule.Wednesday[i] = false
			schedule.Thursday[i] = false
			schedule.Friday[i] = false
		}
	}

func dailyTotal(schedule *Schedule, day string) int {
	total := 0
	switch day {
	case "Monday":
		for i := 0; i < 60; i++ {
			if schedule.Monday[i] {
				total++
			}
		}
	case "Tuesday":
		for i := 0; i < 60; i++ {
			if schedule.Tuesday[i] {
				total++
			}
		}
	case "Wednesday":
		for i := 0; i < 60; i++ {
			if schedule.Wednesday[i] {
				total++
			}
		}
	case "Thursday":
		for i := 0; i < 60; i++ {
			if schedule.Thursday[i] {
				total++
			}
		}
	case "Friday":
		for i := 0; i < 60; i++ {
			if schedule.Friday[i] {
				total++
			}
		}
	}
	return total
}

func weeklyTotal(schedule *Schedule) int {
	total := 0
	total += dailyTotal(schedule, "Monday")
	total += dailyTotal(schedule, "Tuesday")
	total += dailyTotal(schedule, "Wednesday")
	total += dailyTotal(schedule, "Thursday")
	total += dailyTotal(schedule, "Friday")
	return total
}

func convertTimeSlotToHours(slots int) float64 {
	return float64(slots) * 0.16666666
}

func approveSchedule(schedule *Schedule) {
	schedule.Approved = true
}

func toJson(schedule *Schedule) string {
	data, err := json.Marshal(schedule)
	if err != nil {
		fmt.Println("Error converting schedule to JSON:", err)
		return ""
	}
	return string(data)
}

func fromJson(jsonSchedule string) Schedule {
	var schedule Schedule

	err := json.Unmarshal([]byte(jsonSchedule), &schedule)
	if err != nil {
		return Schedule{}
	}

	return schedule
}