package web

import (
	"embed"
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"megen/models"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/browser"
	"github.com/sqweek/dialog"
	"gorm.io/gorm"
)

//go:embed template
var templates embed.FS

var db *gorm.DB

func StartServer(dbNew *gorm.DB) {
	db = dbNew
	t, err := template.ParseFS(templates, "template/*.gohtml")
	if err != nil {
		panic(err)
	}
	r := gin.Default()

	//Startseite
	r.GET("/home", func(c *gin.Context) {
		err = t.ExecuteTemplate(c.Writer, "index.gohtml", nil)
		if err != nil {
			log.Println(err)
		}
	})

	//Orders process
	r.GET("/orders", func(c *gin.Context) {
		var sections models.Sections
		err = db.Debug().Preload("Inputs", "parent_id is null").Preload("Inputs.Children").Where("process = 'orders'").Find(&sections).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		err = t.ExecuteTemplate(c.Writer, "orders.gohtml", sections)
		if err != nil {
			log.Println(err)
		}
	})

	//Delfor process
	r.GET("/delfor", func(c *gin.Context) {
		var sections models.Sections
		err = db.Preload("Inputs", "parent_id is null").Preload("Inputs.Children").Where("process = 'delfor'").Find(&sections).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		err = t.ExecuteTemplate(c.Writer, "orders.gohtml", sections)
		if err != nil {
			log.Println(err)
		}
	})

	//Stockmovement process
	r.GET("/stockmovement", func(c *gin.Context) {
		var sections models.Sections
		err = db.Preload("Inputs", "parent_id is null").Preload("Inputs.Children").Where("process = 'stockmovement'").Find(&sections).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		err = t.ExecuteTemplate(c.Writer, "orders.gohtml", sections)
		if err != nil {
			log.Println(err)
		}
	})

	//Kanban process
	r.GET("/kanban", func(c *gin.Context) {
		var sections models.Sections
		err = db.Preload("Inputs", "parent_id is null").Preload("Inputs.Children").Where("process = 'kanban'").Find(&sections).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		err = t.ExecuteTemplate(c.Writer, "orders.gohtml", sections)
		if err != nil {
			log.Println(err)
		}
	})

	//ASN v2 process
	r.GET("/asnV2", func(c *gin.Context) {
		var sections models.Sections
		err = db.Preload("Inputs", "parent_id is null").Preload("Inputs.Children").Where("process = 'ASN - v2'").Find(&sections).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		err = t.ExecuteTemplate(c.Writer, "orders.gohtml", sections)
		if err != nil {
			log.Println(err)
		}
	})

	//für wenn man es abschickt
	r.POST("/submit", func(c *gin.Context) {
		type submitRequest struct {
			Sections []models.Section
			General  models.General
			Input    models.Input
		}
		var r submitRequest
		err := c.BindJSON(&r)
		if err != nil {
			fmt.Println(err)
			return
		}

		for i := r.General.StartingNumber; i < r.General.NumberOfFiles+r.General.StartingNumber; i++ {
			output := ""
			for _, section := range r.Sections {
				for _, input := range section.Inputs {
					if input.Checked {
						output += input.GetLine(r.General)
					}
				}
			}
			output = postprocess(output, r.General, r.Input)
			filename := fmt.Sprintf(r.General.Directory+"/SO_%s_%s%s_%d.edi", r.General.Process, r.General.UserAbbreviation, time.Now().Format("20060102"), i)
			fmt.Println(filename)
			f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Println("error:", err)
				return
			}
			_, err = f.WriteString(output)
			if err != nil {
				fmt.Println("error:", err)
				return
			}
			f.Close()
			//fmt.Println(output)
		}

	})

	// for Editmode: addSection functionality
	r.GET("/addSection", func(c *gin.Context) {
		type sectionRequest struct {
			Name string `form:"name"`
		}
		var r sectionRequest
		err := c.BindQuery(&r)
		if err != nil {
			fmt.Println(err)
			return
		}

		newSection := models.Section{Name: r.Name}
		err = db.Create(&newSection).Error
		if err != nil {
			fmt.Println(err)
			return
		}
		c.JSON(http.StatusOK, newSection)
	})

	//for Editmode: deleteSection functionality
	r.GET("/deleteSection", func(c *gin.Context) {
		type sectionRequest struct {
			ID uint `form:"id"`
		}
		var r sectionRequest
		err := c.BindQuery(&r)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = db.Delete(&models.Section{}, r.ID).Error
		if err != nil {
			fmt.Println(err)
			return
		}
	})

	//endpoint for directory chooser
	r.GET("/directory", func(c *gin.Context) {
		dir, err := dialog.Directory().Title("Where should the files be saved to?").Browse()
		log.Println(dir)
		if err != nil {
			log.Println(err)
		}
		c.JSON(200, gin.H{"dir": dir})
	})

	port := "localhost:8080"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	browser.OpenURL("http://localhost:8080/home")
	r.Run(port) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

var testcases = []func(edinum string, line string) string{
	// for test case 03: repeat the FTX+AAI segment 12 times;
	// TODO: now only No14 is repeated 12 times. Do the same with No142
	func(edinum, line string) string {
		var tmp = ""
		if edinum == "No14" {
			for i := 0; i < 12; i++ {
				tmp += line + "\n"
			}
		} else {
			return line
		}
		return tmp
	},
	// For test case 04: paste more than 40 characters in IMD+F segment
	// TODO: The same has to be done with No266
	func(edinum, line string) string {
		if edinum == "No126" {
			lineParts := strings.Split(line, ":")
			lineParts[3] = "These are more than forty characters. This should throw an BIS error."
			return strings.Join(lineParts, ":")
		}
		return line
	},
	// For test case 05: exclude the NAD+SU segment
	func(edinum, line string) string {
		if edinum == "No65" {
			line = ""
		}
		return line
	},
	// for test case 09: line is split, org code is replaced, line is merged again
	func(edinum, line string) string {
		if edinum == "No45" {
			lineParts := strings.Split(line, "+")
			lineParts3 := strings.Split(lineParts[3], ":")
			fmt.Println(lineParts)
			fmt.Println(lineParts3)
			lineParts3[0] = "AAAAA"
			lineParts[3] = strings.Join(lineParts3, ":")
			return strings.Join(lineParts, "+")
		}
		return line
	},
}

// postprocess: the lines are splittet in an array, so that at index 0 is the edinumber and at index 1 is the line.
// The Testcases are executet on these arrays and the edinumber is then removed for the final result
func postprocess(input string, g models.General, in models.Input) string {
	lines := strings.Split(input, "\n")

	/*for _, testCase := range testcases {
		for i, line := range lines {
			tmp := strings.Split(line, "\t") // ["NO123", "ABC"]
			if len(tmp) == 2 {
				lines[i] = tmp[0] + "\t" + testCase(tmp[0], tmp[1])
			}
		}
	}*/

	// sort lines by edinumber.
	sort.Slice(lines, func(i, j int) bool {
		tmp1 := strings.Split(lines[i], "\t")
		tmp2 := strings.Split(lines[j], "\t")
		if len(tmp1) == 2 && len(tmp2) == 2 {
			a1, err := strconv.Atoi(tmp1[0][2:]) // "NO123" -> "123"
			if err != nil {
				return false
			}
			a2, err := strconv.Atoi(tmp2[0][2:]) // "NO123" -> "123"
			if err != nil {
				return false
			}
			return a1 < a2
		}
		return false
	})

	// duplicate by number of positions:
	positionStartIndex := 0
	positionEndIndex := 0
	for i := 0; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "No124") {
			positionStartIndex = i
		}
		lParts := strings.Split(lines[i], "\t")
		if ediNumToInt(lParts[0]) > 262 {
			positionEndIndex = i - 1
		}
	}
	posBlock := []string{}
	for i := 0; i < g.NumberOfPositions; i++ {
		for j := positionStartIndex; j < positionEndIndex; j++ {
			posBlock = append(posBlock, strings.ReplaceAll(lines[j], "::POSITIONID::", fmt.Sprintf("%d", i+1)))
		}
	}
	lines = append(append(lines[:positionStartIndex], posBlock...), lines[positionEndIndex:]...)

	// Cutting off the EdiNumbers in front of the lines
	for i, line := range lines {
		tmp := strings.Split(line, "\t")
		if len(tmp) == 2 {
			lines[i] = tmp[1]
		} else {
			lines[i] = line
		}
	}
	return strings.Join(lines, "\n")
}

func ediNumToInt(ediNum string) int {
	if len(ediNum) < 3 {
		return 0
	}
	res, err := strconv.Atoi(ediNum[2:])
	if err != nil {
		panic(err)
	}
	return res
}
