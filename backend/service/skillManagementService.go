package skillservice

import (
	"database/sql"
	"log"
	"net/http"

	skillschemas "github.com/chayutK/skill-management-incubator/backend/schemas"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

var DB *sql.DB
var BindingErrorResponse = gin.H{
	"status":  "error",
	"message": "Binding data error",
}

var BadRequestErrorResponse = gin.H{
	"status":  "error",
	"message": "Input incorrectly",
}
var InternalServerErrorResponse = gin.H{
	"status":  "error",
	"message": "Internal server error",
}

var UpdateErrorResponse = gin.H{
	"status":  "error",
	"message": "not be able to update skill",
}

var UpdateErrorMessage = "not be able to update skill "

var NotFoundErrorResponse = gin.H{
	"status":  "error",
	"message": "Skill not found",
}

var DeleteErrorResponse = gin.H{
	"status":  "error",
	"message": "not be able to delete skill",
}

func HelloWorldHandler(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "Hello, World",
	})
}

func GetAllHandler(ctx *gin.Context) {
	q := "SELECT key,name,description,logo,tags FROM skill"
	rows, err := DB.Query(q)

	if err != nil {
		log.Println("Error while query all skill", err)
		ctx.JSON(http.StatusInternalServerError, InternalServerErrorResponse)
		return
	}

	skills := []skillschemas.Skill{}
	for rows.Next() {
		skill := skillschemas.Skill{}
		err := rows.Scan(&skill.Key, &skill.Name, &skill.Description, &skill.Logo, pq.Array(&skill.Tags))
		if err != nil {
			log.Println("Error not found")
			ctx.JSON(http.StatusNotFound, NotFoundErrorResponse)
			return
		}
		skills = append(skills, skill)

	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   skills,
	})
}

func GetByKeyHandler(ctx *gin.Context) {
	key := ctx.Param("key")

	skill, err := getSkillByKey(key)

	if err != nil {
		log.Println("Error not found", err)
		ctx.JSON(http.StatusNotFound, NotFoundErrorResponse)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   skill,
	})
}

func CreateHandler(ctx *gin.Context) {
	skill := skillschemas.Skill{}
	err := ctx.ShouldBindJSON(&skill)

	if err != nil {
		log.Println("Please fill all the required fields", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Required field not filled",
		})
		return
	}

	q := "INSERT INTO skill (key,name,description,logo,tags) values ($1,$2,$3,$4,$5) RETURNING key,name,description,logo,tags"
	row := DB.QueryRow(q, skill.Key, skill.Name, skill.Description, skill.Logo, pq.Array(skill.Tags))
	result := skillschemas.Skill{}

	err = row.Scan(&result.Key, &result.Name, &result.Description, &result.Logo, pq.Array(&result.Tags))
	if err != nil {
		log.Println("Error while scanning", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Skill already exists",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   result,
	})
}

func UpdateHandler(ctx *gin.Context) {
	key := ctx.Param("key")

	skill := skillschemas.Skill{}
	err := ctx.ShouldBindJSON(&skill)
	skill.Key = key

	if err != nil {
		log.Println("Error while binding data", err)
		ctx.JSON(http.StatusBadRequest, BindingErrorResponse)
	}

	err = updateSkillByKey(skill)

	// q := "UPDATE skill SET name = $1,description = $2, logo=$3,tags=$4 WHERE key=$5 "

	// _, err = DB.Exec(q, skill.Name, skill.Description, skill.Logo, pq.Array(skill.Tags), key)
	// if err != nil {
	// 	log.Println("Error while updating data", err)
	// 	ctx.JSON(http.StatusInternalServerError, UpdateErrorResponse)
	// }

	checkUpdatedSkill, err := getSkillByKey(key)

	if err != nil {
		log.Println("Error while getting data", err)
		ctx.JSON(http.StatusBadRequest, UpdateErrorResponse)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   checkUpdatedSkill,
	})
}

func UpdateNameHandler(ctx *gin.Context) {
	key := ctx.Param("key")

	skill := skillschemas.Skill{}
	err := ctx.BindJSON(&skill)

	if skill.Name == "" {
		log.Println("Name is empty string")
		ctx.JSON(http.StatusBadRequest, BadRequestErrorResponse)
		return
	}

	if err != nil {
		log.Println("Error while binding data", err)
		ctx.JSON(http.StatusInternalServerError, InternalServerErrorResponse)
		return
	}

	oldSkill, err := getSkillByKey(key)

	if err != nil {
		log.Println("Error while getting data", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": UpdateErrorMessage + "name",
		})
		return
	}

	oldSkill.Name = skill.Name

	err = updateSkillByKey(oldSkill)
	if err != nil {
		log.Println("Error while updating data", err)

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": UpdateErrorMessage + "name",
		})
		return
	}

	newSkill, err := getSkillByKey(key)

	if err != nil {
		log.Println("Error while getting data", err)
		ctx.JSON(http.StatusInternalServerError, InternalServerErrorResponse)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   newSkill,
	})

}

func UpdateDescriptionHandler(ctx *gin.Context) {
	key := ctx.Param("key")

	skill := skillschemas.Skill{}
	err := ctx.BindJSON(&skill)

	if skill.Description == "" {
		log.Println("Description is empty string")
		ctx.JSON(http.StatusBadRequest, BadRequestErrorResponse)
		return
	}

	if err != nil {
		log.Println("Error while binding data", err)
		ctx.JSON(http.StatusInternalServerError, InternalServerErrorResponse)
		return
	}

	oldSkill, err := getSkillByKey(key)

	if err != nil {
		log.Println("Error while getting data", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": UpdateErrorMessage + "description",
		})
		return
	}

	oldSkill.Description = skill.Description

	err = updateSkillByKey(oldSkill)
	if err != nil {
		log.Println("Error while updating data", err)

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": UpdateErrorMessage + "description",
		})
		return
	}

	newSkill, err := getSkillByKey(key)

	if err != nil {
		log.Println("Error while getting data", err)
		ctx.JSON(http.StatusInternalServerError, InternalServerErrorResponse)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   newSkill,
	})

}

func UpdateLogoHandler(ctx *gin.Context) {
	key := ctx.Param("key")

	skill := skillschemas.Skill{}
	err := ctx.BindJSON(&skill)

	if skill.Logo == "" {
		log.Println("Logo is empty string")
		ctx.JSON(http.StatusBadRequest, BadRequestErrorResponse)
		return
	}

	if err != nil {
		log.Println("Error while binding data", err)
		ctx.JSON(http.StatusInternalServerError, InternalServerErrorResponse)
		return
	}

	oldSkill, err := getSkillByKey(key)

	if err != nil {
		log.Println("Error while getting data", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": UpdateErrorMessage + "logo",
		})
		return
	}

	oldSkill.Logo = skill.Logo

	err = updateSkillByKey(oldSkill)
	if err != nil {
		log.Println("Error while updating data", err)

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": UpdateErrorMessage + "logo",
		})
		return
	}

	newSkill, err := getSkillByKey(key)

	if err != nil {
		log.Println("Error while getting data", err)
		ctx.JSON(http.StatusInternalServerError, InternalServerErrorResponse)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   newSkill,
	})

}

func UpdateTagsHandler(ctx *gin.Context) {
	key := ctx.Param("key")

	skill := skillschemas.Skill{}
	err := ctx.BindJSON(&skill)

	if len(skill.Tags) == 0 {
		log.Println("Name is empty string")
		ctx.JSON(http.StatusBadRequest, BadRequestErrorResponse)
		return
	}

	if err != nil {
		log.Println("Error while binding data", err)
		ctx.JSON(http.StatusInternalServerError, InternalServerErrorResponse)
		return
	}

	oldSkill, err := getSkillByKey(key)

	if err != nil {
		log.Println("Error while getting data", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": UpdateErrorMessage + "tags",
		})
		return
	}

	oldSkill.Tags = skill.Tags

	err = updateSkillByKey(oldSkill)
	if err != nil {
		log.Println("Error while updating data", err)

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": UpdateErrorMessage + "tags",
		})
		return
	}

	newSkill, err := getSkillByKey(key)

	if err != nil {
		log.Println("Error while getting data", err)
		ctx.JSON(http.StatusInternalServerError, InternalServerErrorResponse)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   newSkill,
	})

}

func DeleteHandler(ctx *gin.Context) {
	key := ctx.Param("key")

	_, err := getSkillByKey(key)
	if err != nil {
		log.Println("Error while getting data", err)
		ctx.JSON(http.StatusBadRequest, DeleteErrorResponse)
		return
	}

	q := "DELETE FROM skill WHERE key=$1"
	_, err = DB.Exec(q, key)

	if err != nil {
		log.Println("Error while deleting data", err)
		ctx.JSON(http.StatusInternalServerError, DeleteErrorResponse)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "sucesss",
		"message": "Skill deleted",
	})

}

// utility function
func getSkillByKey(key string) (skillschemas.Skill, error) {
	q := "SELECT key,name,description,logo,tags FROM skill WHERE key=$1"
	skill := skillschemas.Skill{}
	row := DB.QueryRow(q, key)

	err := row.Scan(&skill.Key, &skill.Name, &skill.Description, &skill.Logo, pq.Array(&skill.Tags))

	if err != nil {
		return skill, err
	}

	return skill, nil
}

func updateSkillByKey(skill skillschemas.Skill) error {
	q := "UPDATE skill SET name = $1,description = $2, logo=$3,tags=$4 WHERE key=$5 "

	_, err := DB.Exec(q, skill.Name, skill.Description, skill.Logo, pq.Array(skill.Tags), skill.Key)
	if err != nil {
		return err
	}
	return nil
}
