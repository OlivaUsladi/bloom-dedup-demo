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
		fmt.Println(*countFlag, *duplicateRatioFlag, *outFlag, *seedFlag, *sourcesFlag)
		model.GenerateEvents(*outFlag, *countFlag, *sourcesFlag, *duplicateRatioFlag, *seedFlag)

	case "run":
		runCmd := flag.NewFlagSet("run", flag.ExitOnError)
		in := runCmd.String("in", "", "входной файл")
		cfg := runCmd.String("config", "", "файл конфигурации")
		outRes := runCmd.String("out", "", "файл результата")
		report := runCmd.String("report", "", "markdown отчёт")
		runCmd.Parse(os.Args[2:])
		fmt.Println(*in, *cfg, *outRes, *report)
	case "bench":
		benchCmd := flag.NewFlagSet("bench", flag.ExitOnError)
		in := benchCmd.String("in", "", "входной файл")
		cfg := benchCmd.String("config", "", "файл конфигурации")

		benchCmd.Parse(os.Args[2:])
		fmt.Println(*in, *cfg)
	default:
		fmt.Println("неизвестная подкоманда:", os.Args[1])
		os.Exit(1)
	}

}
