package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)
type event struct {
	Acronym      string `json:"acronym"`
	ActivityCode string `json:"activityCode"`
	Artermin     string `json:"artermin"`
	CourseCode   string `json:"courseCode"`
	CourseName   struct {
		NameNob string `json:"nameNob"`
		NameNno string `json:"nameNno"`
		NameEng string `json:"nameEng"`
	} `json:"courseName"`
	Name     string        `json:"name"`
	Termnr   int           `json:"termnr"`
	Title    string        `json:"title"`
	TpID     string        `json:"tpId"`
	Disiplin []interface{} `json:"disiplin"`
	From     int64         `json:"from"`
	Staff    []interface{} `json:"staff"`
	Rooms    []struct {
		ID       string `json:"id"`
		Building string `json:"building"`
		Room     string `json:"room"`
		URL      string `json:"url"`
	} `json:"rooms"`
	Status           string   `json:"status"`
	StudyProgramKeys []string `json:"studyProgramKeys"`
	Summary          string   `json:"summary"`
	To               int64    `json:"to"`
	Week             int      `json:"week"`
	SelectedProgram  string   `json:"selectedProgram"`
}

type events struct {
	Schedules []struct {
		Acronym      string `json:"acronym"`
		ActivityCode string `json:"activityCode"`
		Artermin     string `json:"artermin"`
		CourseCode   string `json:"courseCode"`
		CourseName   struct {
			NameNob string `json:"nameNob"`
			NameNno string `json:"nameNno"`
			NameEng string `json:"nameEng"`
		} `json:"courseName"`
		Name     string        `json:"name"`
		Termnr   int           `json:"termnr"`
		Title    string        `json:"title"`
		TpID     string        `json:"tpId"`
		Disiplin []interface{} `json:"disiplin"`
		From     int64         `json:"from"`
		Staff    []interface{} `json:"staff"`
		Rooms    []struct {
			ID       string `json:"id"`
			Building string `json:"building"`
			Room     string `json:"room"`
			URL      string `json:"url"`
		} `json:"rooms"`
		Status           string   `json:"status"`
		StudyProgramKeys []string `json:"studyProgramKeys"`
		Summary          string   `json:"summary"`
		To               int64    `json:"to"`
		Week             int      `json:"week"`
		SelectedProgram  string   `json:"selectedProgram"`
	} `json:"schedules"`
}



func main() {
	router := gin.Default()
	fileserver := http.FileServer(http.Dir("./static"))
	router.GET("/events/:courseCode", func(c *gin.Context) {
		courseCode := c.Param("courseCode")
		retrieveEvents(c, courseCode)
	})

	router.GET("eventsThisWeek/:courseCode", func(c *gin.Context) {
		courseCode := c.Param("courseCode")
		getEventsThisWeek(c, courseCode)
	})

	router.GET("/eventsByCourseCodes/:stringOfCourseCodes", func(c *gin.Context) {
		stringOfCourseCodes := c.Param("stringOfCourseCodes")
		log.Println(stringOfCourseCodes)
		courseCodes := strings.Split(stringOfCourseCodes, ",")
		log.Println(courseCodes)


		getListOfCoursesByCourseCodes(c, courseCodes)
	})

	router.GET("/eventsByCourseCodesByWeek/:week/:stringOfCourseCodes", func(c *gin.Context) {
		stringOfCourseCodes := c.Param("stringOfCourseCodes")
		week := c.Param("week")

		//convert week to int
		weekInt, err := strconv.Atoi(week)
		if err != nil {
			log.Println(err)
		}

		

		log.Println(stringOfCourseCodes)
		courseCodes := strings.Split(stringOfCourseCodes, ",")
		log.Println(courseCodes)

		getListOfCoursesByCourseCodesByWeek(c, courseCodes, weekInt)
	})

	


	
	// serve the html, css and js files from the static folder
	router.GET("/", func(c *gin.Context) {
		fileserver.ServeHTTP(c.Writer, c.Request)
	})
	router.Static("/static", "./static")




	router.Run("localhost:3000")
}


func retrieveEvents(c *gin.Context, courseCode string) {
	//get events from external api
	// https://www.ntnu.no/web/studier/emner?p_p_id=coursedetailsportlet_WAR_courselistportlet&p_p_lifecycle=2&p_p_state=normal&p_p_mode=view&p_p_resource_id=schedules&p_p_cacheability=cacheLevelPage&_coursedetailsportlet_WAR_courselistportlet_courseCode=IDATA2302&year=2022&version=1
	// replace courseCode with the courseCode from the request
	resp, err := fetchAllEvents(courseCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var events events
	err = json.Unmarshal(body, &events)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, events)
}

func fetchAllEvents(courseCode string) (*http.Response, error) {
	url := "https://www.ntnu.no/web/studier/emner?p_p_id=coursedetailsportlet_WAR_courselistportlet&p_p_lifecycle=2&p_p_state=normal&p_p_mode=view&p_p_resource_id=schedules&p_p_cacheability=cacheLevelPage&_coursedetailsportlet_WAR_courselistportlet_courseCode=" + courseCode + "&year=2022&version=1"
	resp, err := http.Get(url)
	return resp, err
}


func getEventsThisWeek(c *gin.Context, courseCode string) {
	
	//use existing function to get events
	resp, err := fetchAllEvents(courseCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var events events
	err = json.Unmarshal(body, &events)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// filter events to only include events this week
	// get current week
	thisWeek := Week(time.Now())

	// add events to array if they are this week
	var eventsThisWeek []event
	for _, event := range events.Schedules {
		if event.Week == thisWeek {
			eventsThisWeek = append(eventsThisWeek, event)
		}
	}
	c.JSON(http.StatusOK, eventsThisWeek)

}

func getListOfCoursesByCourseCodes(c *gin.Context, courseCodes []string) {

	var listOfEvents []events
	for _, courseCode := range courseCodes {
		resp, err := fetchAllEvents(courseCode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var events events
		err = json.Unmarshal(body, &events)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		listOfEvents = append(listOfEvents, events)
	}
	//merge all schedules into one array
	var mergedSchedules []event
	for _, events := range listOfEvents {
		for _, schedule := range events.Schedules {
			mergedSchedules = append(mergedSchedules, schedule)
		}
	}

	c.JSON(http.StatusOK, mergedSchedules)

}

func getListOfCoursesByCourseCodesByWeek(c *gin.Context, courseCodes []string, week int) {

	var listOfEvents []events
	for _, courseCode := range courseCodes {
		resp, err := fetchAllEvents(courseCode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var events events
		err = json.Unmarshal(body, &events)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		listOfEvents = append(listOfEvents, events)
	}
	//merge all schedules into one array
	var mergedSchedules []event
	for _, events := range listOfEvents {
		for _, schedule := range events.Schedules {
			mergedSchedules = append(mergedSchedules, schedule)
		}
	}

	//filter by week
	var eventsThisWeek []event
	for _, event := range mergedSchedules {
		if event.Week == week {
			eventsThisWeek = append(eventsThisWeek, event)
		}
	}

	c.JSON(http.StatusOK, eventsThisWeek)

}





func Week(now time.Time) int {
	_, week := now.ISOWeek()
	return week
}