package main

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"
)

const (
	Z1 = 1000
	Z2 = 1000
)

var b1 = "A"
var b2 = "B"

type Loja struct {
	Ngjarja1 float64
	Ngjarja2 float64
}

var lista = map[string]map[string]Loja{}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}
func GenTestData() {
	rand.Seed(time.Now().UTC().UnixNano())
	B := 15
	L := 10
	for i := 0; i <= B; i++ {
		l := map[string]Loja{}
		for j := 0; j <= L; j++ {
			r := float64(random(1, 10)) / float64(random(1, 10))
			if r <= 1 {
				r *= float64(random(3, 10))
			}
			r1 := float64(random(1, 10)) / float64(random(1, 10))
			if r1 <= 1 {
				r1 *= float64(random(3, 10))
			}
			l[fmt.Sprintf("L%d", j)] = Loja{r, r1}
		}
		lista[fmt.Sprintf("B%d", i)] = l
	}
}
func TestMod(t *testing.T) {
	GenTestData()
	t1 := time.Now()
	final_verdict := ""
	diff := 0.0
	for k1, b1 := range lista {
		for k2, b2 := range lista {
			if k1 != k2 {
				for lk, l1 := range b1 {
					x11 := l1.Ngjarja1
					y12 := b2[lk].Ngjarja2

					x12 := b2[lk].Ngjarja1
					y11 := l1.Ngjarja2
					for z1 := 1; z1 <= Z1; z1++ {
						for z2 := 1; z2 <= Z2; z2++ {
							if (float64(z1)*x11 > float64(z1+z2)) && (float64(z2)*y12 > float64(z1+z2)) {
								//fmt.Println(fmt.Sprintf("Loja -> %s, B-> %s Ngjarja 1,%s  Ngjarja 2, Z1 -> %d, Z2 -> %d",lk,k1,k2,z1,z2))
								v := math.Sqrt((float64(z1)*x11 - float64(z1+z2)) * ((float64(z2) * y12) - float64(z1+z2)))
								v = v * (1 / (math.Abs(float64(z1)*x11-float64(z2)*y12) + 1))
								if v > diff {
									final_verdict = ""
									diff = v
									final_verdict += fmt.Sprintln(fmt.Sprintf("Loja -> %s, BS -> %s Ngjarja 1 - %f, BS -> %s  Ngjarja 2 - %f, Z1 -> %d, Z2 -> %d, ", lk, k1,x11, k2, y12, z1, z2))
									final_verdict += fmt.Sprintln("If Match 1 ->", float64(z1)*x11-float64(z1+z2))
									final_verdict += fmt.Sprintln("If Match 2 ->", float64(z2)*y12-float64(z1+z2))
								}
							} else if float64(z1)*x11 <= float64(z1+z2) {
								break
							}
						}
					}
					for z1 := 1; z1 <= Z1; z1++ {
						for z2 := 1; z2 <= Z2; z2++ {
							if (float64(z1)*y11 > float64(z1+z2)) && (float64(z2)*x12 > float64(z1+z2)) {
								//fmt.Println(fmt.Sprintf("Loja -> %s, B-> %s Ngjarja 2,%s  Ngjarja 1, Z1 -> %d, Z2 -> %d",lk,k1,k2,z1,z2))
								v := math.Sqrt((float64(z1)*y11 - float64(z1+z2)) * ((float64(z2) * x12) - float64(z1+z2)))
								v = v * (1 / (math.Abs(float64(z1)*y11-float64(z2)*x12) + 1))
								if v > diff {
									final_verdict = ""
									diff = v
									final_verdict += fmt.Sprintln(fmt.Sprintf("Loja -> %s, BS -> %s Ngjarja 2 - %f, BS -> %s  Ngjarja 1 - %f, Z1 -> %d, Z2 -> %d, ", lk, k1,y11, k2, x12, z1, z2))
									final_verdict += fmt.Sprintln("If Match 1 ->", float64(z1)*y11-float64(z1+z2))
									final_verdict += fmt.Sprintln("If Match 2 ->", float64(z2)*x12-float64(z1+z2))
								}
							} else if float64(z1)*y11 <= float64(z1+z2) {
								break
							}
						}
					}
				}
			}
		}
	}
	fmt.Println(time.Since(t1))
	fmt.Println(final_verdict)
}

//func TestName(t *testing.T) {
//	t1 := time.Now()
//	diff := 0.0
//	final_vedict := ""
//	for i := 1; i <= Z1; i++ {
//		for i2 := 1; i2 <= Z2; i2++ {
//			for l, v := range lista {
//				for _, b := range v {
//					for _, b1 := range v {
//						if b.Id != b1.Id {
//							if float64(i)*b.Ngjarja1 > float64(i+i2) && float64(i2)*b1.Ngjarja2 > float64(i+i2) {
//								v := math.Sqrt((float64(i)*b.Ngjarja1 - float64(i+i2)) * ((float64(i2) * b1.Ngjarja2) - float64(i+i2)))
//								v = v * (1 / (math.Abs(float64(i)*b.Ngjarja1-float64(i2)*b1.Ngjarja2) + 1))
//								if v > diff {
//									final_vedict = ""
//									diff = v
//									final_vedict += fmt.Sprintln(fmt.Sprintf("Match ---> 1. Bastore: %s, Loja: %s, Vlera: %d, Ngjarja 1", b.Id, l, i))
//									final_vedict += fmt.Sprintln(fmt.Sprintf("Match ---> 2. Bastore: %s, Loja: %s, Vlera: %d, Ngjarja 2", b1.Id, l, i2))
//									final_vedict += fmt.Sprintln("If Match 1 ->", float64(i)*b.Ngjarja1-float64(i+i2))
//									final_vedict += fmt.Sprintln("If Match 2 ->", float64(i2)*b1.Ngjarja2-float64(i+i2))
//								}
//							}
//						}
//					}
//				}
//			}
//		}
//	}
//	fmt.Println(time.Since(t1))
//	fmt.Println()
//	fmt.Println("Final Vedict")
//	fmt.Println(final_vedict)
//}
