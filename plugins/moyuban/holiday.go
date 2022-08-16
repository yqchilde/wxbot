package moyuban

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/tidwall/gjson"
)

var (
	NextHoliday   = make(map[string]struct{}, 3)
	ContentHeader = "摸鱼人！工作再累，一定不要忘记摸鱼哦！有事没事起身去茶水间，去厕所，去廊道走走别老在工位上坐着，钱是老板的，但命是自己的。"
	ContentFooter = "上班是帮老板赚钱，摸鱼是赚老板的钱！最后，祝愿天下所有摸鱼人，都能愉快地度过每一天..."
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func DailyLifeNotes(date string, caller ...int) (string, error) {
	if len(caller) == 0 {
		caller = append(caller, 1)
	}
	_, callerFile, _, _ := runtime.Caller(caller[0])
	dayFile, err := ioutil.ReadFile(filepath.Dir(callerFile) + "/holiday.json")
	if err != nil {
		panic(err)
	}

	var (
		currentTime    = time.Now().Local()
		toWeekend      int
		holidayBalance string
		loafOnJobNotes string
		weekendNotes   string
		holidayNotes   string
		vacationNotes  string
		dailyLifeNotes string
		workInWeekend  bool
	)

	if len(date) > 0 {
		parse, err := time.Parse("2006-01-02 15:04:05", date)
		if err != nil {
			return "", err
		}
		currentTime = parse
	}
	toWeekend = int(6 - currentTime.Weekday())
	holiday := gjson.Get(string(dayFile), "holiday")
	holiday.ForEach(func(key, val gjson.Result) bool {
		holidayDate, _ := time.Parse("2006-01-02", val.Get("date").String())
		festivalName := val.Get("name").String()

		if _, ok := NextHoliday[festivalName]; ok {
			return true
		}
		if val.Get("holiday").Bool() && currentTime.Sub(holidayDate).Seconds() < 0 {
			days := int(holidayDate.Sub(currentTime).Hours()/24 + 1)
			NextHoliday[val.Get("name").String()] = struct{}{}
			if currentTime.Hour() < 12 {
				holidayBalance += fmt.Sprintf("距离【%s】还有：%d天\n", val.Get("name").String(), days)
			} else {
				holidayBalance += fmt.Sprintf("距离【%s】还有：%d天\n", val.Get("name").String(), days)
			}
		}
		if !val.Get("holiday").Bool() {
			for i := 1; i <= toWeekend; i++ {
				addDate := currentTime.AddDate(0, 0, i)
				if val.Get("date").String() == addDate.Format("2006-01-02") && (addDate.Weekday() == time.Saturday || addDate.Weekday() == time.Sunday) {
					toWeekend++
					workInWeekend = true
				}
			}
		}
		return true
	})
	NextHoliday = make(map[string]struct{})

	COVID19Duration := currentTime.Sub(time.Date(2019, 12, 16, 0, 0, 0, 0, time.Local)).Hours() / 24
	COVID19Data := fmt.Sprintf("自新冠疫情爆发以来已经过了%d天了，疫情防控形势再度严峻，注意及时配合防疫政策，做好自我防护！", int(COVID19Duration))
	workInWeekendStr := ""
	if workInWeekend {
		workInWeekendStr = "(受调休影响)"
	}
	if currentTime.Hour() < 12 {
		loafOnJobNotes = fmt.Sprintf("【摸鱼办】提醒您：\n%d月%d日上午好，%s\n%s\n距离【周末】还有：%d天%s\n%s%s", currentTime.Month(), currentTime.Day(), ContentHeader, COVID19Data, toWeekend, workInWeekendStr, holidayBalance, ContentFooter)
		weekendNotes = fmt.Sprintf("【摸鱼办】提醒您：\n%d月%d日上午好，\n今天是周末，多出去走走吧\n%s", currentTime.Month(), currentTime.Day(), holidayBalance)
		holidayNotes = fmt.Sprintf("【摸鱼办】提醒您：\n%d月%d日上午好，\n今天是假期，祝你玩的愉快\n%s", currentTime.Month(), currentTime.Day(), holidayBalance)
		vacationNotes = fmt.Sprintf("【摸鱼办】提醒您：\n%d月%d日上午好，\n今天是节假日调休\n%s", currentTime.Month(), currentTime.Day(), holidayBalance)
	} else {
		loafOnJobNotes = fmt.Sprintf("【摸鱼办】提醒您：\n%d月%d日下午好，%s\n%s\n距离【周末】还有：%d天%s\n%s%s", currentTime.Month(), currentTime.Day(), ContentHeader, COVID19Data, toWeekend, workInWeekendStr, holidayBalance, ContentFooter)
		weekendNotes = fmt.Sprintf("【摸鱼办】提醒您：\n%d月%d日下午好，\n今天是周末，多出去走走吧\n%s", currentTime.Month(), currentTime.Day(), holidayBalance)
		holidayNotes = fmt.Sprintf("【摸鱼办】提醒您：\n%d月%d日下午好，\n今天是假期，祝你玩的愉快\n%s", currentTime.Month(), currentTime.Day(), holidayBalance)
		vacationNotes = fmt.Sprintf("【摸鱼办】提醒您：\n%d月%d日下午好，\n今天是节假日调休\n%s", currentTime.Month(), currentTime.Day(), holidayBalance)
	}

	// 判断文件是否存在
	var infoBytes []byte
	exists, err := PathExists(fmt.Sprintf("info-%s.json", currentTime.Format("20060102")))
	if err != nil {
		return "", err
	}
	if exists {
		readFile, err := ioutil.ReadFile(fmt.Sprintf("info-%s.json", currentTime.Format("20060102")))
		if err != nil {
			return "", err
		}
		infoBytes = readFile
	} else {
		resp, err := http.Get(fmt.Sprintf("https://timor.tech/api/holiday/info/%s", currentTime.Format("2006-01-02")))
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		infoBytes = body
		ioutil.WriteFile(fmt.Sprintf("info-%s.json", currentTime.Format("20060102")), body, 0644)
	}

	result := gjson.Get(string(infoBytes), "type.type")
	switch result.Int() {
	case 0:
		dailyLifeNotes = loafOnJobNotes
	case 1:
		dailyLifeNotes = weekendNotes
	case 2:
		dailyLifeNotes = holidayNotes
	case 3:
		dailyLifeNotes = vacationNotes
	}
	return dailyLifeNotes, nil
}
