package logic

import (
	"fmt"
	"math/big"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func Home(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func GetPayslip(c *gin.Context) {
	grossStr := c.PostForm("gross")
	healthStr := c.PostForm("health")
	contribStr := c.PostForm("contrib")
	housingStr := c.PostForm("housing")

	gross, err := strconv.ParseFloat(grossStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid gross value"})
		return
	}

	health := healthStr == "true"
	contrib := contribStr == "true"
	housing := housingStr == "true"

	// fmt.Printf("gross: %f\n", gross)
	// fmt.Printf("health: %t, %T\n", health, health)
	// fmt.Printf("contrib: %t, %T\n", contrib, contrib)
	// fmt.Printf("housing: %t, %T\n", housing, housing)

	grossBig := big.NewFloat(gross)
	grade := NewGrade(grossBig, health, contrib, housing)
	// fmt.Printf("Grade: %t, %T\n", grade, grade)

	monthlyPay := grade.GetNetPay()
	payee := grade.MonthlyPayee()
	healthIns := grade.GetHealthEmpl()
	housings := grade.GetHousingAmount()
	emp_pension := grade.GetMonthlyEmployeePension()
	pension := grade.GetMonthlyPension()

	response := gin.H{
		"payslip2":     fmt.Sprintf("%.2f", bigFloatToFloat64(monthlyPay)),
		"payee2":       fmt.Sprintf("%.2f", bigFloatToFloat64(payee)),
		"health2":      fmt.Sprintf("%.2f", bigFloatToFloat64(new(big.Float).Quo(healthIns, big.NewFloat(12)))),
		"housing2":     fmt.Sprintf("%.2f", bigFloatToFloat64(new(big.Float).Quo(housings, big.NewFloat(12)))),
		"pension2":     fmt.Sprintf("%.2f", bigFloatToFloat64(pension)),
		"emp_pension2": fmt.Sprintf("%.2f", bigFloatToFloat64(emp_pension)),
		"payslip":      fmt.Sprintf("%s", FormatWithThousandSeparator(monthlyPay)),
		"payee":        fmt.Sprintf("%s", FormatWithThousandSeparator(payee)),
		"health":       fmt.Sprintf("%s", FormatWithThousandSeparator(healthIns)),
		"housing":      fmt.Sprintf("%s", FormatWithThousandSeparator(housings)),
		"emp_pension":  fmt.Sprintf("%s", FormatWithThousandSeparator(emp_pension)),
		"pension":      fmt.Sprintf("%s", FormatWithThousandSeparator(pension)),
	}

	c.JSON(http.StatusOK, response)
}

func bigFloatToFloat64(f *big.Float) float64 {
	result, _ := f.Float64()
	return result
}
