package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Input struct {
	Values []float64 `json:"values"`
}

// Функція для обчислення параметрів палива
func calculateTask1(Hp, Cp, Sp, Np, Op, Wp, Ap float64) string {
	// Обчислення коефіцієнтів переходу
	Krs := 100 / (100 - Wp) 
	Krg := 100 / (100 - Wp - Ap)

	// Обчислення складу сухої маси
	Hc := Hp * Krs
	Cc := Cp * Krs
	Sc := Sp * Krs
	Nc := Np * Krs
	Oc := Op * Krs
	Ac := Ap * Krs

	// Обчислення складу горючої маси
	Hg := Hp * Krg
	Cg := Cp * Krg
	Sg := Sp * Krg
	Ng := Np * Krg
	Og := Op * Krg

	// Обчислення нижчої теплоти згоряння
	Qrn := (339*Cp + 1030*Hp - 108.8*(Op-Sp) - 25*Wp) / 1000
	Qsn := (Qrn + 0.025*Wp) * 100 / (100 - Wp)
	Qhn := (Qrn + 0.025*Wp) * 100 / (100 - Wp - Ap)

	// Вивід результату
	return fmt.Sprintf("Компонентний склад - H^P: %.2f%%, C^P: %.2f%%, S^P: %.2f%%, N^P: %.2f%%, O^P: %.2f%%, W^P: %.2f%%, A^P: %.2f%%\n"+
		"\nКоеф. переходу від робочої до сухої маси - %.2f\n"+
		"\nКоеф. переходу від робочої до горючої маси - %.2f\n"+
		"\nСклад сухої маси палива - H^C: %.2f%%, C^C: %.2f%%, S^C: %.2f%%, N^C: %.2f%%, O^C: %.2f%%, A^C: %.2f%%\n"+
		"\nСклад горючої маси палива - H^G: %.2f%%, C^G: %.2f%%, S^G: %.2f%%, N^G: %.2f%%, O^G: %.2f%%\n"+
		"Н. теплота згоряння для робочої маси - %.2f МДж/кг\n"+
		"Н. теплота згоряння для сухої маси - %.2f МДж/кг\n"+
		"Н. теплота згоряння для горючої маси - %.2f МДж/кг",
		Hp, Cp, Sp, Np, Op, Wp, Ap,
		Krs, Krg, Hc, Cc, Sc, Nc, Oc, Ac,
		Hg, Cg, Sg, Ng, Og,
		Qrn, Qsn, Qhn,
	)
}

// Функція для обчислення характеристик мазуту
func calculateTask2(Hg, Cg, Og, Sg, Qdaf, Wg, Ag, Vg float64) string {
	// Обчислення складу робочої маси
	Cp := Cg * (100 - Wg - Ag) / 100
	Hp := Hg * (100 - Wg - Ag) / 100
	Op := Og * (100 - Wg - Ag) / 100
	Sp := Sg * (100 - Wg*0.1 - Ag*0.1) / 100
	Ap := Ag * (100 - Wg) / 100
	Vp := Vg * (100 - Wg) / 100

	// Обчислення нижчої теплоти згоряння
	Qp := Qdaf*((100-Wg-Ap)/100) - 0.025*Wg

	// Формування результату
	return fmt.Sprintf(
		"Склад горючої маси мазуту: H^Г: %.2f%%, C^Г: %.2f%%, S^Г: %.2f%%, O^Г: %.2f%%, V^Г: %.2f мг/кг, W^Г: %.2f%%, A^Г: %.2f%%;\n\n"+
		"Склад робочої маси мазуту: H^Р: %.2f%%, C^Р: %.2f%%, S^Р: %.2f%%, O^Р: %.2f%%, V^Р: %.2f мг/кг, A^Р: %.2f%%;\n\n"+
		"Нижча теплота згоряння мазуту на робочу масу для робочої маси: %.2f МДж/кг.",
		Hg, Cg, Sg, Og, Vg, Wg, Ag,
		Hp, Cp, Sp, Op, Vp, Ap, Qp,
	)
}

// Обробник запитів для першого калькулятора
func calculator1Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var input Input
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if len(input.Values) != 7 {
		http.Error(w, "Invalid number of inputs", http.StatusBadRequest)
		return
	}

	result := calculateTask1(input.Values[0], input.Values[1], input.Values[2], input.Values[3],
		input.Values[4], input.Values[5], input.Values[6])

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"result": result})
}

// Обробник запитів для другого калькулятора
func calculator2Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var input Input
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if len(input.Values) != 8 {
		http.Error(w, "Invalid number of inputs", http.StatusBadRequest)
		return
	}

	result := calculateTask2(input.Values[0], input.Values[1], input.Values[2], input.Values[3],
		input.Values[4], input.Values[5], input.Values[6], input.Values[7])

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"result": result})
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/api/calculator1", calculator1Handler)
	http.HandleFunc("/api/calculator2", calculator2Handler)

	fmt.Println("Server running at http://localhost:8081")
	http.ListenAndServe(":8081", nil)
}
