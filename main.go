package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
)

type Mes struct {
	AporteRF         float64 `json:"aporte_rf"`
	AporteFIIs       float64 `json:"aporte_fiis"`
	ValorBrutoRF     float64 `json:"valor_bruto_rf"`
	ValorLiquidoRF   float64 `json:"valor_liquido_rf"`
	ValorLiquidoFIIs float64 `json:"valor_liquido_fiis"`
}

type Ano map[string]Mes

type Dados struct {
	Anos map[string]Ano `json:"anos"`
}

const arquivo = "dados.json"

func carregarDados() Dados {
	file, err := os.ReadFile(arquivo)
	if err != nil {
		return Dados{Anos: make(map[string]Ano)}
	}

	var dados Dados
	err = json.Unmarshal(file, &dados)
	if err != nil {
		fmt.Println("Erro ao carregar dados:", err)
		return Dados{Anos: make(map[string]Ano)}
	}
	return dados
}

func salvarDados(dados Dados) {
	bytes, err := json.MarshalIndent(dados, "", "  ")
	if err != nil {
		fmt.Println("Erro ao salvar dados:", err)
		return
	}
	os.WriteFile(arquivo, bytes, 0644)
}

func menu() {
	dados := carregarDados()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\n--- MENU PRINCIPAL ---")
		fmt.Println("1. Ver resumo completo")
		fmt.Println("2. Adicionar/editar mês")
		fmt.Println("3. Sair")
		fmt.Print("Escolha uma opção: ")
		scanner.Scan()
		opcao := scanner.Text()

		switch opcao {
		case "1":
			mostrarResumo(dados)
		case "2":
			adicionarOuEditarMes(&dados, scanner)
			salvarDados(dados)
		case "3":
			fmt.Println("Saindo...")
			return
		default:
			fmt.Println("Opção inválida!")
		}
	}
}

func nomeMes(m string) string {
	nomes := map[string]string{
		"01": "Janeiro", "02": "Fevereiro", "03": "Março",
		"04": "Abril", "05": "Maio", "06": "Junho",
		"07": "Julho", "08": "Agosto", "09": "Setembro",
		"10": "Outubro", "11": "Novembro", "12": "Dezembro",
	}
	return nomes[m]
}

func mostrarResumo(dados Dados) {
	fmt.Println("\n📌 Resumo dos aportes e saldos mensais")

	totalAporte := 0.0
	valorBrutoRFAcumulado := 0.0
	valorLiquidoRFAcumulado := 0.0
	valorLiquidoFIIsAcumulado := 0.0

	ultimoAno := ""
	ultimoMes := ""
	var mesAtual Mes

	for ano, meses := range dados.Anos {
		for mes := range meses {
			if ano > ultimoAno || (ano == ultimoAno && mes > ultimoMes) {
				ultimoAno = ano
				ultimoMes = mes
				mesAtual = meses[mes]
			}
		}
	}

	fmt.Printf("\n🗓️  MÊS ATUAL: %s/%s\n", nomeMes(ultimoMes), ultimoAno)
	fmt.Printf("Aporte RF: R$ %.2f | Aporte FIIs: R$ %.2f\n", mesAtual.AporteRF, mesAtual.AporteFIIs)
	fmt.Printf("Valor Bruto RF: R$ %.2f\n", mesAtual.ValorBrutoRF)
	fmt.Printf("Valor Líquido RF: R$ %.2f | Valor Líquido FIIs: R$ %.2f\n", mesAtual.ValorLiquidoRF, mesAtual.ValorLiquidoFIIs)
	fmt.Printf("Lucro Bruto RF: R$ %.2f\n", mesAtual.ValorBrutoRF-mesAtual.AporteRF)
	fmt.Printf("Lucro Líquido Total: R$ %.2f\n", mesAtual.ValorLiquidoRF+mesAtual.ValorLiquidoFIIs-mesAtual.AporteRF-mesAtual.AporteFIIs)

	fmt.Println("\n| Mês      | Aporte Total | Aporte RF | Aporte FIIs | Valor Bruto RF | Valor Líquido RF | Valor Líquido FIIs | Lucro Bruto Ac. | Lucro Líquido Ac. |")
	fmt.Println("|----------|--------------|-----------|-------------|----------------|------------------|--------------------|------------------|--------------------|")

	anos := ordenarChaves(dados.Anos)
	aporteRFSoFar := 0.0
	aporteFIIsSoFar := 0.0

	for _, ano := range anos {
		meses := ordenarChaves(dados.Anos[ano])
		for _, mes := range meses {
			m := dados.Anos[ano][mes]

			aporteRFSoFar += m.AporteRF
			aporteFIIsSoFar += m.AporteFIIs
			valorBrutoRFAcumulado = m.ValorBrutoRF
			valorLiquidoRFAcumulado = m.ValorLiquidoRF
			valorLiquidoFIIsAcumulado = m.ValorLiquidoFIIs

			lucroBrutoAcumulado := m.ValorBrutoRF - aporteRFSoFar
			lucroLiquidoAcumulado := (m.ValorLiquidoRF + m.ValorLiquidoFIIs) - (aporteRFSoFar + aporteFIIsSoFar)

			fmt.Printf("| %-8s | R$ %10.2f | R$ %7.2f | R$ %9.2f | R$ %14.2f | R$ %16.2f | R$ %18.2f | R$ %16.2f | R$ %18.2f |\n",
				nomeMes(mes), m.AporteRF+m.AporteFIIs, m.AporteRF, m.AporteFIIs,
				m.ValorBrutoRF, m.ValorLiquidoRF, m.ValorLiquidoFIIs, lucroBrutoAcumulado, lucroLiquidoAcumulado)
		}
	}

	totalAporte = aporteRFSoFar + aporteFIIsSoFar
	totalLucroBrutoRF := valorBrutoRFAcumulado - aporteRFSoFar
	totalLucroLiquido := (valorLiquidoRFAcumulado + valorLiquidoFIIsAcumulado) - totalAporte

	fmt.Printf("\nTotal aportado: R$ %.2f\n", totalAporte)
	fmt.Printf("Valor bruto final (RF): R$ %.2f\n", valorBrutoRFAcumulado)
	fmt.Printf("Valor líquido final (RF): R$ %.2f\n", valorLiquidoRFAcumulado)
	fmt.Printf("Valor líquido final (FIIs): R$ %.2f\n", valorLiquidoFIIsAcumulado)
	fmt.Printf("Lucro bruto total (RF): R$ %.2f\n", totalLucroBrutoRF)
	fmt.Printf("Lucro líquido total: R$ %.2f\n", totalLucroLiquido)
}

func ordenarChaves[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func adicionarOuEditarMes(dados *Dados, scanner *bufio.Scanner) {
	fmt.Print("Digite o ano (ex: 2025): ")
	scanner.Scan()
	ano := scanner.Text()

	fmt.Print("Digite o mês (ex: 05): ")
	scanner.Scan()
	mes := scanner.Text()

	fmt.Print("Digite o aporte na Renda Fixa: R$ ")
	scanner.Scan()
	aporteRF, _ := strconv.ParseFloat(scanner.Text(), 64)

	fmt.Print("Digite o aporte em FIIs: R$ ")
	scanner.Scan()
	aporteFIIs, _ := strconv.ParseFloat(scanner.Text(), 64)

	fmt.Print("Digite o valor bruto da Renda Fixa: R$ ")
	scanner.Scan()
	valorBrutoRF, _ := strconv.ParseFloat(scanner.Text(), 64)

	fmt.Print("Digite o valor líquido da Renda Fixa: R$ ")
	scanner.Scan()
	valorLiquidoRF, _ := strconv.ParseFloat(scanner.Text(), 64)

	fmt.Print("Digite o valor líquido dos FIIs: R$ ")
	scanner.Scan()
	valorLiquidoFIIs, _ := strconv.ParseFloat(scanner.Text(), 64)

	if dados.Anos[ano] == nil {
		dados.Anos[ano] = make(Ano)
	}

	dados.Anos[ano][mes] = Mes{
		AporteRF:         aporteRF,
		AporteFIIs:       aporteFIIs,
		ValorBrutoRF:     valorBrutoRF,
		ValorLiquidoRF:   valorLiquidoRF,
		ValorLiquidoFIIs: valorLiquidoFIIs,
	}

	fmt.Println("Dados salvos com sucesso!")
}

func main() {
	menu()
}
