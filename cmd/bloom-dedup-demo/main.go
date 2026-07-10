package main

import (
	"bloom-dedup-demo/internal/model"
	"flag"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("нужна подкоманда: generate, run или bench")
		os.Exit(1)
	}
	fmt.Println(os.Args[1])
	switch os.Args[1] {
	case "generate":
		genCmd := flag.NewFlagSet("generate", flag.ExitOnError)
		countFlag := genCmd.Int("count", 0, "общее число записей")
		duplicateRatioFlag := genCmd.Float64("duplicate-ratio", 0.0, "доля дублей")
		outFlag := genCmd.String("out", "", "файл вывода")
		seedFlag := genCmd.Int64("seed", 0, "Seed генератора")
		sourcesFlag := genCmd.Int("sources", 1, "число источников")
		genCmd.Parse(os.Args[2:])
		if *countFlag <= 0 {
			fmt.Fprintln(os.Stderr, "count должен быть больше 0")
			os.Exit(1)
		}
		if *duplicateRatioFlag < 0 || *duplicateRatioFlag > 0.9 {
			fmt.Fprintln(os.Stderr, "duplicate-ratio от 0 до 0.9")
			os.Exit(1)
		}
		if *outFlag == "" {
			fmt.Fprintln(os.Stderr, "путь не может быть пустым")
			os.Exit(1)
		}
		if *sourcesFlag <= 0 || *sourcesFlag > 100 {
			fmt.Fprintln(os.Stderr, "источники должны быть от 1 до 100")
			os.Exit(1)
		}
		fmt.Println(*countFlag, *duplicateRatioFlag, *outFlag, *seedFlag, *sourcesFlag)
		model.GenerateEvents(*outFlag, *countFlag, *sourcesFlag, *duplicateRatioFlag, *seedFlag)

	case "run":
		runCmd := flag.NewFlagSet("run", flag.ExitOnError)
		in := runCmd.String("in", "", "входной файл")
		cfg := runCmd.String("config", "", "файл конфигурации")
		outRes := runCmd.String("out", "", "файл результата")
		report := runCmd.String("report", "", "markdown отчёт")
		sourcesBoolFlag := runCmd.Bool("fls", true, "флаг пропуска событий с невалидными источником и датой (true - пропуск)")
		runCmd.Parse(os.Args[2:])
		if *in == "" {
			fmt.Fprintln(os.Stderr, "путь файла событий не может быть пустым")
			os.Exit(1)
		}
		if *cfg == "" {
			fmt.Fprintln(os.Stderr, "путь файла конфигурации не может быть пустым")
			os.Exit(1)
		}
		if *outRes == "" {
			fmt.Fprintln(os.Stderr, " путь выходного файла не может быть пустым")
			os.Exit(1)
		}
		//if *report != "" {
		//	fmt.Fprintln(os.Stderr, "путь отчёта не должен быть пустым")
		//	os.Exit(1)
		//}
		fmt.Println(*in, *cfg, *outRes, *report, *sourcesBoolFlag)
	case "bench":
		benchCmd := flag.NewFlagSet("bench", flag.ExitOnError)
		in := benchCmd.String("in", "", "входной файл")
		cfg := benchCmd.String("config", "", "файл конфигурации")
		sourcesBoolFlag := benchCmd.Bool("fls", true, "флаг пропуска событий с невалидными источником и датой (true - пропуск)")

		benchCmd.Parse(os.Args[2:])
		if *in == "" {
			fmt.Fprintln(os.Stderr, "путь файла событий не может быть пустым")
			os.Exit(1)
		}
		if *cfg == "" {
			fmt.Fprintln(os.Stderr, "путь файла конфигурации не должен быть пустым")
			os.Exit(1)
		}
		fmt.Println(*in, *cfg, *sourcesBoolFlag)
	default:
		fmt.Println("неизвестная подкоманда:", os.Args[1])
		os.Exit(1)
	}

}
