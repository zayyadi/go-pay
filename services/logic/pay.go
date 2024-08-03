package logic

import (
	"fmt"
	"math/big"
	"strings"
	"time"
)

type Variable struct {
	HEALTH_EMPL_PERC        *big.Float
	HEALTH_EMPLYR_PERC      *big.Float
	THREE_SIXTY             *big.Float
	TWO_HUND                *big.Float
	NON_TAXABLE_VARIABLES   *big.Float
	FIRST_TAXABLE_VARIABLES *big.Float
	SIX_K                   *big.Float
	FIVE_K                  *big.Float
	ONE_ONE_M               *big.Float
	ONE_SIX_M               *big.Float
	THREE_TWO_M             *big.Float
	SIX_FOUR_M              *big.Float
	SEVENTY_K               *big.Float
	HOUSING_PERC            *big.Float
}

var varConfig = Variable{
	HEALTH_EMPL_PERC:        big.NewFloat(5),   // Example value
	HEALTH_EMPLYR_PERC:      big.NewFloat(7.5), // Example value
	THREE_SIXTY:             big.NewFloat(360000),
	TWO_HUND:                big.NewFloat(200000),
	NON_TAXABLE_VARIABLES:   big.NewFloat(30000),
	FIRST_TAXABLE_VARIABLES: big.NewFloat(300000),
	SIX_K:                   big.NewFloat(600000),
	FIVE_K:                  big.NewFloat(500000),
	ONE_ONE_M:               big.NewFloat(1100000),
	ONE_SIX_M:               big.NewFloat(1600000),
	THREE_TWO_M:             big.NewFloat(3200000),
	SIX_FOUR_M:              big.NewFloat(6400000),
	SEVENTY_K:               big.NewFloat(70000),
	HOUSING_PERC:            big.NewFloat(10),
}

type Grade struct {
	Gross     *big.Float
	WaterFee  *big.Float
	HealthIns bool
	Contrib   bool
	Housing   bool
}

func NewGrade(gross *big.Float, healthIns, contrib, housing bool) *Grade {
	return &Grade{
		Gross:     gross,
		WaterFee:  big.NewFloat(150),
		HealthIns: healthIns,
		Contrib:   contrib,
		Housing:   housing,
	}
}

func (g *Grade) GetAnnualGross() *big.Float {
	annualGross := new(big.Float).Mul(g.Gross, big.NewFloat(12))
	return annualGross
}

func (g *Grade) getBasic() *big.Float {
	basic := new(big.Float).Mul(g.GetAnnualGross(), big.NewFloat(0.4))
	return basic
}

func (g *Grade) getTransport() *big.Float {
	transport := new(big.Float).Mul(g.GetAnnualGross(), big.NewFloat(0.1))
	return transport
}

func (g *Grade) getHousing() *big.Float {
	housing := new(big.Float).Mul(g.GetAnnualGross(), big.NewFloat(0.1))
	return housing
}

func (g *Grade) GetHealthEmpl() *big.Float {
	if g.HealthIns {
		return new(big.Float).Mul(g.Gross, new(big.Float).Quo(varConfig.HEALTH_EMPL_PERC, big.NewFloat(100)))
	}
	return big.NewFloat(0)
}

func (g *Grade) GetHealthEmplyr() *big.Float {
	if g.HealthIns {
		return new(big.Float).Mul(g.Gross, new(big.Float).Quo(varConfig.HEALTH_EMPLYR_PERC, big.NewFloat(100)))
	}
	return big.NewFloat(0)
}

func (g *Grade) GetHealthIns() *big.Float {
	return new(big.Float).Add(g.GetHealthEmpl(), g.GetHealthEmplyr())
}

func (g *Grade) GetBHT() *big.Float {
	bht := new(big.Float).Add(g.getBasic(), g.getHousing())
	bht.Add(bht, g.getTransport())
	return bht
}

func (g *Grade) GetPensionEmployees() *big.Float {
	annualGross := g.GetAnnualGross()
	if annualGross.Cmp(varConfig.THREE_SIXTY) <= 0 {
		return big.NewFloat(0)
	}
	if g.Contrib {
		return new(big.Float).Mul(annualGross, big.NewFloat(0.08))
	}
	return new(big.Float).Mul(annualGross, big.NewFloat(0.18))
}

func (g *Grade) GetPensionEmployer() *big.Float {
	annualGross := g.GetAnnualGross()
	if annualGross.Cmp(varConfig.THREE_SIXTY) <= 0 {
		return big.NewFloat(0)
	}
	if g.Contrib {
		return new(big.Float).Mul(annualGross, big.NewFloat(0.10))
	}
	return big.NewFloat(0)
}

func (g *Grade) GetMonthlyPension() *big.Float {
	return new(big.Float).Quo(g.GetTotalPension(), big.NewFloat(12))
}

func (g *Grade) GetMonthlyEmployeePension() *big.Float {
	return new(big.Float).Quo(g.GetPensionEmployees(), big.NewFloat(12))
}

func (g *Grade) GetTotalPension() *big.Float {
	return new(big.Float).Add(g.GetPensionEmployees(), g.GetPensionEmployer())
}

func (g *Grade) PensionLogic() *big.Float {
	annualGross := g.GetAnnualGross()
	if annualGross.Cmp(varConfig.THREE_SIXTY) <= 0 {
		return big.NewFloat(0)
	}
	return g.GetPensionEmployees()
}

func (g *Grade) GetHousingAmount() *big.Float {
	if g.Housing {
		return new(big.Float).Mul(g.Gross, new(big.Float).Quo(varConfig.HOUSING_PERC, big.NewFloat(100)))
	}
	return big.NewFloat(0)
}

func (g *Grade) GetNsitf() *big.Float {
	return new(big.Float).Mul(g.Gross, big.NewFloat(0.01))
}

func (g *Grade) GetGrossIncome() *big.Float {
	grossIncome := new(big.Float).Sub(g.GetAnnualGross(), g.PensionLogic())
	return grossIncome
}

func (g *Grade) TwentyPercents() *big.Float {
	return new(big.Float).Mul(g.GetGrossIncome(), big.NewFloat(0.2))
}

func (g *Grade) GetConsolidated() *big.Float {
	annualGross := g.GetAnnualGross()
	if new(big.Float).Mul(annualGross, big.NewFloat(0.01)).Cmp(varConfig.TWO_HUND) > 0 {
		return new(big.Float).Mul(annualGross, big.NewFloat(0.01))
	}
	return varConfig.TWO_HUND
}

func (g *Grade) GetConsolidatedRelief() *big.Float {
	return new(big.Float).Add(g.GetConsolidated(), g.TwentyPercents())
}

func (g *Grade) GetTaxableIncome() *big.Float {
	taxableIncome := new(big.Float).Sub(g.GetAnnualGross(), g.GetConsolidatedRelief())
	taxableIncome.Sub(taxableIncome, g.GetPensionEmployees())
	if taxableIncome.Cmp(big.NewFloat(0)) <= 0 {
		return big.NewFloat(0)
	}
	return taxableIncome
}

func (g *Grade) FirstTaxable() *big.Float {
	if g.GetTaxableIncome().Cmp(varConfig.NON_TAXABLE_VARIABLES) <= 0 {
		return big.NewFloat(0)
	}
	return big.NewFloat(0)
}

func (g *Grade) SecondTaxable() *big.Float {
	taxableIncome := g.GetTaxableIncome()
	if taxableIncome.Cmp(varConfig.FIRST_TAXABLE_VARIABLES) < 0 {
		return new(big.Float).Mul(taxableIncome, big.NewFloat(0.07))
	} else if taxableIncome.Cmp(varConfig.FIRST_TAXABLE_VARIABLES) >= 0 {
		return new(big.Float).Mul(varConfig.FIRST_TAXABLE_VARIABLES, big.NewFloat(0.07))
	}
	return big.NewFloat(0)
}

func (g *Grade) ThirdTaxable() *big.Float {
	taxableIncome := new(big.Float).Sub(g.GetTaxableIncome(), varConfig.FIRST_TAXABLE_VARIABLES)
	if taxableIncome.Cmp(big.NewFloat(0)) > 0 && taxableIncome.Cmp(varConfig.FIRST_TAXABLE_VARIABLES) <= 0 {
		return new(big.Float).Mul(taxableIncome, big.NewFloat(0.11))
	} else if taxableIncome.Cmp(varConfig.FIRST_TAXABLE_VARIABLES) > 0 {
		return new(big.Float).Mul(varConfig.FIRST_TAXABLE_VARIABLES, big.NewFloat(0.11))
	}
	return big.NewFloat(0)
}

func (g *Grade) FourthTaxable() *big.Float {
	taxableIncome := new(big.Float).Sub(g.GetTaxableIncome(), varConfig.SIX_K)
	if taxableIncome.Cmp(varConfig.FIVE_K) >= 0 {
		return new(big.Float).Mul(varConfig.FIVE_K, big.NewFloat(0.15))
	} else if taxableIncome.Cmp(big.NewFloat(0)) > 0 && taxableIncome.Cmp(varConfig.FIVE_K) <= 0 {
		return new(big.Float).Mul(taxableIncome, big.NewFloat(0.15))
	}
	return big.NewFloat(0)
}

func (g *Grade) FifthTaxable() *big.Float {
	taxableIncome := new(big.Float).Sub(g.GetTaxableIncome(), varConfig.ONE_ONE_M)
	if taxableIncome.Cmp(varConfig.FIVE_K) >= 0 {
		return new(big.Float).Mul(varConfig.FIVE_K, big.NewFloat(0.19))
	} else if taxableIncome.Cmp(big.NewFloat(0)) > 0 && taxableIncome.Cmp(varConfig.FIVE_K) < 0 {
		return new(big.Float).Mul(taxableIncome, big.NewFloat(0.19))
	}
	return big.NewFloat(0)
}

func (g *Grade) SixthTaxable() *big.Float {
	taxableIncome := new(big.Float).Sub(g.GetTaxableIncome(), varConfig.ONE_SIX_M)
	if taxableIncome.Cmp(varConfig.ONE_SIX_M) >= 0 {
		return new(big.Float).Mul(varConfig.ONE_SIX_M, big.NewFloat(0.21))
	} else if taxableIncome.Cmp(big.NewFloat(0)) > 0 && taxableIncome.Cmp(varConfig.ONE_SIX_M) < 0 {
		return new(big.Float).Mul(taxableIncome, big.NewFloat(0.21))
	}
	return big.NewFloat(0)
}

func (g *Grade) SeventhTaxable() *big.Float {
	taxableIncome := new(big.Float).Sub(g.GetTaxableIncome(), varConfig.THREE_TWO_M)
	if taxableIncome.Cmp(varConfig.THREE_TWO_M) > 0 {
		return new(big.Float).Mul(taxableIncome, big.NewFloat(0.24))
	} else if taxableIncome.Cmp(big.NewFloat(0)) <= 0 {
		return big.NewFloat(0)
	}
	return big.NewFloat(0)
}

func (g *Grade) PayeeLogic() *big.Float {
	taxableIncome := g.GetTaxableIncome()
	if taxableIncome.Cmp(varConfig.NON_TAXABLE_VARIABLES) <= 0 {
		return big.NewFloat(0)
	} else if taxableIncome.Cmp(varConfig.FIRST_TAXABLE_VARIABLES) <= 0 {
		return g.SecondTaxable()
	} else if taxableIncome.Cmp(varConfig.FIRST_TAXABLE_VARIABLES) >= 0 && taxableIncome.Cmp(varConfig.SIX_K) < 0 {
		return new(big.Float).Add(g.SecondTaxable(), g.ThirdTaxable())
	} else if taxableIncome.Cmp(varConfig.FIRST_TAXABLE_VARIABLES) >= 0 && taxableIncome.Cmp(varConfig.SIX_K) >= 0 && taxableIncome.Cmp(varConfig.ONE_ONE_M) < 0 {
		return new(big.Float).Add(new(big.Float).Add(g.SecondTaxable(), g.ThirdTaxable()), g.FourthTaxable())
	} else if taxableIncome.Cmp(varConfig.ONE_ONE_M) >= 0 && taxableIncome.Cmp(varConfig.ONE_SIX_M) < 0 {
		return new(big.Float).Add(new(big.Float).Add(new(big.Float).Add(g.SecondTaxable(), g.ThirdTaxable()), g.FourthTaxable()), g.FifthTaxable())
	} else if taxableIncome.Cmp(varConfig.ONE_SIX_M) >= 0 && taxableIncome.Cmp(varConfig.THREE_TWO_M) < 0 {
		return new(big.Float).Add(new(big.Float).Add(new(big.Float).Add(new(big.Float).Add(g.SecondTaxable(), g.ThirdTaxable()), g.FourthTaxable()), g.FifthTaxable()), g.SixthTaxable())
	} else if taxableIncome.Cmp(varConfig.SIX_FOUR_M) >= 0 {
		return new(big.Float).Add(new(big.Float).Add(new(big.Float).Add(new(big.Float).Add(new(big.Float).Add(g.SecondTaxable(), g.ThirdTaxable()), g.FourthTaxable()), g.FifthTaxable()), g.SixthTaxable()), g.SeventhTaxable())
	}
	return big.NewFloat(0)
}

func (g *Grade) GetWaterFee() *big.Float {
	if g.Gross.Cmp(varConfig.SEVENTY_K) >= 0 {
		return g.WaterFee
	} else {
		return big.NewFloat(200)
	}
}

func (g *Grade) MonthlyPayee() *big.Float {
	return new(big.Float).Quo(g.PayeeLogic(), big.NewFloat(12))
}

func (g *Grade) GetNetPay() *big.Float {
	netPay := new(big.Float).Sub(g.Gross, g.MonthlyPayee())
	netPay.Sub(netPay, new(big.Float).Quo(g.PensionLogic(), big.NewFloat(12)))
	netPay.Sub(netPay, g.GetHealthEmpl())
	netPay.Sub(netPay, g.GetWaterFee())
	netPay.Sub(netPay, g.GetHousingAmount())
	return netPay
}

func (g *Grade) String() string {
	return fmt.Sprintf("Your salary for the month of %s is %s and your net pay is %.2f and your PAYEE for the month is %.2f",
		time.Now().Month().String(), g.Gross.Text('f', 2), g.GetNetPay(), new(big.Float).Quo(g.PayeeLogic(), big.NewFloat(12)))
}

func FormatWithThousandSeparator(f *big.Float) string {
	// Convert the big.Float to a string
	str := f.Text('f', 2) // Fixed-point notation with 2 decimal places

	// Split the integer and fractional parts
	parts := strings.Split(str, ".")
	integerPart := parts[0]
	fractionalPart := ""
	if len(parts) > 1 {
		fractionalPart = parts[1]
	}

	// Insert thousand separators into the integer part
	var result strings.Builder
	for i, digit := range integerPart {
		if i > 0 && (len(integerPart)-i)%3 == 0 {
			result.WriteRune(',')
		}
		result.WriteRune(digit)
	}

	// Append the fractional part, if any
	if fractionalPart != "" {
		result.WriteRune('.')
		result.WriteString(fractionalPart)
	}

	return result.String()
}

func main() {
	var salary float64
	fmt.Print("Enter your salary: ")
	fmt.Scan(&salary)

	gross := big.NewFloat(salary)
	emp := NewGrade(gross, false, true, false)

	fmt.Println(emp.String())
	fmt.Printf("Annual Income: ₦%.2f\n", emp.GetAnnualGross())
	fmt.Printf("Gross: ₦%.2f\n", emp.Gross)
	fmt.Printf("Water: ₦%.2f\n", emp.GetWaterFee())
	fmt.Printf("Housing: ₦%.2f\n", emp.GetHousingAmount())
	fmt.Printf("Cons relief: ₦%.2f\n", emp.GetConsolidatedRelief())
	formatCons := FormatWithThousandSeparator(emp.GetConsolidated())
	fmt.Printf("Cons: ₦%s\n", formatCons)
	fmt.Printf("First taxable: %.2f\n", emp.FirstTaxable())
	fmt.Printf("Second taxable: %.2f\n", emp.SecondTaxable())
	fmt.Printf("Third taxable: %.2f\n", emp.ThirdTaxable())
	fmt.Printf("Fourth taxable: %.2f\n", emp.FourthTaxable())
	fmt.Printf("Fifth taxable: %.2f\n", emp.FifthTaxable())
	fmt.Printf("Sixth taxable: %.2f\n", emp.SixthTaxable())
	fmt.Printf("Seventh taxable: %.2f\n", emp.SeventhTaxable())
	fmt.Printf("Taxable Income: %.2f\n", emp.GetTaxableIncome())
	fmt.Printf("Health: ₦%.2f\n", new(big.Float).Quo(emp.GetHealthEmpl(), big.NewFloat(12)))
	fmt.Printf("Employee Pension Contribution: ₦%.2f\n", new(big.Float).Quo(emp.GetPensionEmployees(), big.NewFloat(12)))
	fmt.Printf("Employee Pension logic: ₦%.2f\n", new(big.Float).Quo(emp.PensionLogic(), big.NewFloat(12)))
	fmt.Printf("Employer's pension contribution: ₦%.2f\n", emp.GetPensionEmployer())
	fmt.Printf("Employee Gross income: ₦%.2f\n", emp.GetGrossIncome())
	fmt.Printf("Net pay for the year: ₦%.2f\n", emp.GetNetPay())
}
