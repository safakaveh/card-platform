package jalali

import (
	"time"
)

const defNumber = 492269

type JalaliConverter struct {
	days          int64
	yearMonthDay  [3]int
	hourMinuteSec [3]int
}

// Constructor با زمان میلادی
func NewJalaliFromMillis(timeMillis int64) *JalaliConverter {
	days := (timeMillis / (1000 * 60 * 60 * 24)) + defNumber
	j := &JalaliConverter{
		days: days,
	}
	j.yearMonthDay = j.getYearMonthDay()
	return j
}

// Constructor با سال، ماه، روز، ساعت، دقیقه، ثانیه
func NewJalali(year, month, day, hour, minute, second int) *JalaliConverter {
	j := &JalaliConverter{
		hourMinuteSec: [3]int{hour - 4, minute - 30, second},
	}
	j.days = j.getDays(year, month, day)
	j.yearMonthDay = j.getYearMonthDay()
	return j
}

// Constructor با سال، ماه، روز
func NewJalaliYMD(year, month, day int) *JalaliConverter {
	return NewJalali(year, month, day, 0, 0, 0)
}

// Constructor بدون ورودی (زمان فعلی)
func NewJalaliNow() *JalaliConverter {
	return NewJalaliFromMillis(time.Now().UnixMilli())
}

// Getters
func (j *JalaliConverter) Year() int {
	return j.yearMonthDay[0]
}

func (j *JalaliConverter) Month() int {
	return j.yearMonthDay[1]
}

func (j *JalaliConverter) Day() int {
	return j.yearMonthDay[2]
}

// تبدیل به میلی‌ثانیه
func (j *JalaliConverter) TimeMillis() int64 {
	oneSecondMs := int64(1000)
	oneMinuteMs := 60 * oneSecondMs
	oneHourMs := 60 * oneMinuteMs
	oneDayMs := 24 * oneHourMs

	return (j.days-defNumber)*oneDayMs +
		int64(j.hourMinuteSec[0])*oneHourMs +
		int64(j.hourMinuteSec[1])*oneMinuteMs +
		int64(j.hourMinuteSec[2])*oneSecondMs
}

// بررسی سال کبیسه جلالی
func (j *JalaliConverter) IsJalaliLeapYear(year int) bool {
	r := year % 33
	return r == 1 || r == 5 || r == 9 || r == 13 || r == 17 || r == 22 || r == 26 || r == 30
}

// طول ماه جلالی
func (j *JalaliConverter) getJalaliMonthLength(year, month int) int {
	if month <= 6 {
		return 31
	}
	if month <= 11 {
		return 30
	}
	if j.IsJalaliLeapYear(year) {
		return 30
	}
	return 29
}

// محاسبه سال، ماه، روز از تعداد روزها
func (j *JalaliConverter) getYearMonthDay() [3]int {
	var result [3]int
	last := j.days
	year := 1

	for last > 366 {
		if j.IsJalaliLeapYear(year) {
			last -= 366
		} else {
			last -= 365
		}
		year++
	}
	if last == 366 && !j.IsJalaliLeapYear(year) {
		last -= 365
		year++
	}
	result[0] = int(year)

	month := 1
	for last > int64(j.getJalaliMonthLength(int(year), month)) {
		last -= int64(j.getJalaliMonthLength(int(year), month))
		month++
	}
	result[1] = month
	result[2] = int(last)

	return result
}

// محاسبه تعداد روزها از تاریخ
func (j *JalaliConverter) getDays(year, month, day int) int64 {
	cntLeapDay := 0
	for i := 0; i < year; i++ {
		if j.IsJalaliLeapYear(i) {
			cntLeapDay++
		}
	}

	if year == 0 {
		year = 1
	}

	dy := int64((year-1)*365 + cntLeapDay)
	for i := 1; i < month; i++ {
		dy += int64(j.getJalaliMonthLength(year, i))
	}
	dy += int64(day)
	return dy
}
