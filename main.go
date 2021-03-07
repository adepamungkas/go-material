package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const startupMessage = `starting...`
var db *gorm.DB
func init() {
	//open a db connection
	var err error
	//db, err = gorm.Open("mysql", "doadmin:d77gjs15qccjf34k@(dbmysql-do-user-8830210-0.b.db.ondigitalocean.com)/defaultdb?charset=utf8&parseTime=True&loc=Local")
	//user:password@(localhost)/dbname?charset=utf8&parseTime=True&loc=Local
	//mysql://doadmin:k8bxpa54neic0gd5@db-mysql-sgp1-09976-do-user-8830210-0.b.db.ondigitalocean.com:25060/defaultdb?ssl-mode=REQUIRED
	//db, err = gorm.Open("mysql", "root:@(localhost:3306)/misdb?charset=utf8&parseTime=True&loc=Local")
	db, err = gorm.Open("mysql", "doadmin:k8bxpa54neic0gd5@10.104.0.0/20:25060/defaultdb?ssl-mode=REQUIRED")
	//db, err = gorm.Open("postgres", "host=http://127.0.0.1:19153 port=5432 user=postgres dbname=materialdb password=12345")

	if err != nil {
		panic(err)
	}

	//Migrate the schema
	db.AutoMigrate(&materialModel{})
}

func main() {

	router := gin.Default()
	router.Use(cors.Default())

	v1 := router.Group("/api/v1/materials")
	{
		v1.POST("/", createMaterial)
		v1.GET("/", getAllMaterial)
		v1.GET("/:ID", getMaterialById)
		v1.PUT("/:ID", updateMaterial)
		v1.DELETE("/:ID", deleteMaterial)
	}
	router.Run()

	//http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	fmt.Fprintf(w, "Hello! you've requested %s\n", r.URL.Path)
	//})
	//
	//http.HandleFunc("/cached", func(w http.ResponseWriter, r *http.Request) {
	//	maxAgeParams, ok := r.URL.Query()["max-age"]
	//	if ok && len(maxAgeParams) > 0 {
	//		maxAge, _ := strconv.Atoi(maxAgeParams[0])
	//		w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", maxAge))
	//	}
	//	requestID := uuid.Must(uuid.NewV4())
	//	fmt.Fprintf(w, requestID.String())
	//})
	//
	//http.HandleFunc("/headers", func(w http.ResponseWriter, r *http.Request) {
	//	keys, ok := r.URL.Query()["key"]
	//	if ok && len(keys) > 0 {
	//		fmt.Fprintf(w, r.Header.Get(keys[0]))
	//		return
	//	}
	//	headers := []string{}
	//	for key, values := range r.Header {
	//		headers = append(headers, fmt.Sprintf("%s=%s", key, strings.Join(values, ",")))
	//	}
	//	fmt.Fprintf(w, strings.Join(headers, "\n"))
	//})
	//
	//http.HandleFunc("/env", func(w http.ResponseWriter, r *http.Request) {
	//	keys, ok := r.URL.Query()["key"]
	//	if ok && len(keys) > 0 {
	//		fmt.Fprintf(w, os.Getenv(keys[0]))
	//		return
	//	}
	//	envs := []string{}
	//	for _, env := range os.Environ() {
	//		envs = append(envs, env)
	//	}
	//	fmt.Fprintf(w, strings.Join(envs, "\n"))
	//})
	//
	//http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
	//	codeParams, ok := r.URL.Query()["code"]
	//	if ok && len(codeParams) > 0 {
	//		statusCode, _ := strconv.Atoi(codeParams[0])
	//		if statusCode >= 200 && statusCode < 600 {
	//			w.WriteHeader(statusCode)
	//		}
	//	}
	//	requestID := uuid.Must(uuid.NewV4())
	//	fmt.Fprintf(w, requestID.String())
	//})

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	for _, encodedRoute := range strings.Split(os.Getenv("ROUTES"), ",") {
		if encodedRoute == "" {
			continue
		}
		pathAndBody := strings.SplitN(encodedRoute, "=", 2)
		path, body := pathAndBody[0], pathAndBody[1]
		http.HandleFunc("/"+path, func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, body)
		})
	}

	bindAddr := fmt.Sprintf(":%s", port)
	lines := strings.Split(startupMessage, "\n")
	fmt.Println()
	for _, line := range lines {
		fmt.Println(line)
	}
	fmt.Println()
	fmt.Printf("==> Server listening at %s 🚀\n", bindAddr)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		panic(err)
	}
}

type (
	// materialModel describes a materialModel type
	materialModel struct {
		gorm.Model
		Trademark     string `json:"Trademark"`
		IsBroken int    `json:"IsBroken"`
		Color string `json:"Color"`
		Date time.Time `json:"Date"`
		Description string `json:"Description"`
		InputBy string `json:"InputBy"`
		Name string `json:"Name"`
		Size int `json:"Size"`
		Type string `json:"Type"`
		Vendor string `json:"Vendor"`


	}

	// materialViewModel represents a formatted material
	materialViewModel struct {
		ID        uint   `json:"ID"`
		Trademark     string `json:"Trademark"`
		IsBroken bool   `json:"IsBroken"`
		Color string `json:"Color"`
		Date time.Time `json:"Date" time_format:"unix"`
		Description string `json:"Description"`
		InputBy string `json:"InputBy"`
		Name string `json:"Name"`
		Size int `json:"Size"`
		Type string `json:"Type"`
		Vendor string `json:"Vendor"`
	}
)
// createMaterial add a new material
func createMaterial(c *gin.Context) {

	isBroken, _ := strconv.Atoi(c.PostForm("IsBroken"))
	size, _ :=strconv.Atoi(c.PostForm("Size"))

	date, _ := time.Parse(time.RFC822Z, c.PostForm("Date"))
	material := materialModel{
		Trademark: c.PostForm("Trademark"),
		IsBroken: isBroken,
		Color:c.PostForm("Color"),
		Date: date,
		Description: c.PostForm("Description"),
		InputBy: c.PostForm("InputBy"),
		Name: c.PostForm("Name"),
		Type: c.PostForm("Type"),
		Vendor: c.PostForm("Vendor"),
		Size: size,

	}
	db.Save(&material)
	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Material item created successfully!", "resourceId": material.ID})
}

// getAllMaterial fetch all materials
func getAllMaterial(c *gin.Context) {

	var materials []materialModel
	var materialInfo []materialViewModel

	db.Find(&materials)

	if len(materials) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No material found!"})
		return
	}

	//transforms the material for building a good response
	for _, item := range materials {
		isBroken := false
		if item.IsBroken == 1 {
			isBroken = true
		} else {
			isBroken = false
		}
		materialInfo = append(materialInfo,
			materialViewModel{
				ID: item.ID,
				Trademark: item.Trademark,
				IsBroken: isBroken,
				Type: item.Type,
				Name: item.Name,
				Vendor: item.Vendor,
				InputBy: item.InputBy,
				Description: item.Description,
				Color: item.Color,
				Size: item.Size,
				Date: item.Date,
			})
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": materialInfo})
}

// getMaterialById fetch a single material
func getMaterialById(c *gin.Context) {
	var material materialModel
	materialID := c.Param("ID")

	db.First(&material, materialID)

	if material.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No material found!"})
		return
	}

	isBroken := false
	if material.IsBroken == 1 {
		isBroken = true
	} else {
		isBroken = false
	}

	materialInfo := materialViewModel{
		ID: material.ID,
		Trademark: material.Trademark,
		IsBroken: isBroken,
		Date: material.Date,
		Size: material.Size,
		Color: material.Color,
		InputBy: material.InputBy,
		Description: material.Description,
		Name: material.Name,
		Vendor: material.Vendor,
		Type: material.Type,
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": materialInfo})
}



// updateMaterial update a material
func updateMaterial(c *gin.Context) {
	var material materialModel
	materialID := c.Param("ID")

	db.First(&material, materialID)

	if material.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No material found!"})
		return
	}

	db.Model(&material).Update(
		"Trademark", c.PostForm("Trademark"),
		"Vendor",c.PostForm("Vendor"),
		"Color",c.PostForm("Color"),
		"Color",c.PostForm("Type"),
		"Description",c.PostForm("Description"),
		"Name",c.PostForm("Name"),
		"isBroken",c.PostForm("isBroken"),
		"Size",c.PostForm("Size"),
	)
	isBroken, _ := strconv.Atoi(c.PostForm("IsBroken"))
	db.Model(&material).Update("IsBroken", isBroken)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Material updated successfully!"})
}

// deleteMaterial remove a material
func deleteMaterial(c *gin.Context) {
	var material materialModel
	materialID := c.Param("ID")

	db.First(&material, materialID)

	if material.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No material found!"})
		return
	}

	db.Delete(&material)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Material deleted successfully!"})
}