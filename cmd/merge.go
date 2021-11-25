package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

var (
	mergeCmd = &cobra.Command{
		Use:     "merge DIRECTORY -o OUTPUT",
		Aliases: []string{"m"},
		Short:   "Merge Videos or Files",
		Long:    "Merge Videos or Files",
		Example: `vip-video-downloader merge /to/path -o output.ts`,
		Run:     mergeRun,
	}

	output     string
	sortType   string
	sortBySeq  bool
	sortByTime bool
	remove     bool
)

func init() {
	mergeCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "Output file")
	mergeCmd.PersistentFlags().StringVarP(&sortType, "sort-type", "t", "asc", "Sort type[asc or desc]")
	mergeCmd.PersistentFlags().BoolVarP(&sortBySeq, "sort-by-seq", "S", false, "Files are in DIRECTORY sorts by seq[0,1,2,3,4,...]")
	mergeCmd.PersistentFlags().BoolVarP(&sortByTime, "sort-by-time", "T", false, "Files are in DIRECTORY sorts by time")
	mergeCmd.PersistentFlags().BoolVarP(&remove, "remove", "r", false, "Remove DIRECTORY after merge done")
	mergeCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "V", false, "Print verbose log")
}

func mergeRun(_ *cobra.Command, args []string) {
	if len(args) <= 0 {
		panic("error: require DIRECTORY")
	}
	if output == "" {
		panic("error: require output file")
	}
	dir := args[0]
	{
		if stat, err := os.Stat(dir); err != nil {
			panic(err)
		} else if !stat.IsDir() {
			panic("error: require DIRECTORY not FILE")
		}
	}
	dirFiles, err := ioutil.ReadDir(dir)
	{
		if err != nil {
			panic(err)
		}
		if len(dirFiles) == 0 {
			panic("error: DIRECTORY is empty")
		}
	}
	seqNums := make([]int, 0)
	seqNames := make([]string, 0)
	times := make([]int, 0)
	fileMap := make(map[string]string, 0)
	timeFileMap := make(map[int]string, 0)
	for _, f := range dirFiles {
		name := f.Name()
		times = append(times, int(f.ModTime().UnixMilli()))
		timeFileMap[int(f.ModTime().UnixMilli())] = f.Name()
		if extIndex := strings.LastIndexByte(name, '.'); extIndex != -1 {
			name = f.Name()[:extIndex]
		}
		if seq, err2 := strconv.Atoi(name); err2 != nil {
			seqNames = append(seqNames, name)
			fileMap[name] = f.Name()
		} else {
			seqNums = append(seqNums, seq)
			fileMap[fmt.Sprintf("%d", seq)] = f.Name()
		}
	}
	orderedFiles := make([]string, 0)
	if sortBySeq {
		sort.Ints(seqNums)
		sort.Strings(seqNames)
		if sortType == "desc" {
			sort.Reverse(sort.IntSlice(seqNums))
			sort.Reverse(sort.StringSlice(seqNames))
		}
		for _, s := range seqNums {
			if p, have := fileMap[fmt.Sprintf("%d", s)]; have {
				orderedFiles = append(orderedFiles, p)
			}
		}
		for _, s := range seqNames {
			if p, have := fileMap[s]; have {
				orderedFiles = append(orderedFiles, p)
			}
		}
	} else if sortByTime {
		sort.Ints(times)
		if sortType == "desc" {
			sort.Reverse(sort.IntSlice(times))
		}
		for _, t := range times {
			if p, have := timeFileMap[t]; have {
				orderedFiles = append(orderedFiles, p)
			}
		}
	}
	mergedFile, err := os.OpenFile(output, os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0700)
	if err != nil {
		panic(err)
	}
	defer func() { _ = mergedFile.Close() }()
	for _, f := range orderedFiles {
		if openFile, err := os.OpenFile(filepath.Join(dir, f), os.O_RDONLY, 0700); err != nil {
			panic(err)
		} else {
			if _, err := io.Copy(mergedFile, openFile); err != nil {
				panic(err)
			}
			if verbose {
				fmt.Println(fmt.Sprintf("[merge] %s", filepath.Join(dir, f)))
			}
			_ = openFile.Close()
		}
	}
	if remove {
		_ = os.RemoveAll(dir)
	}
	fmt.Println(fmt.Sprintf("[success] %s", output))
}
