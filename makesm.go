/*
 * I B A U K   -   S C O R E M A S T E R
 *
 * I make standalone, distributable, installations of ScoreMaster
 * ready to be burned to CD/DVD/USB
 *
 * I am written for readability rather than efficiency, please keep me that way.
 *
 *
 * Copyright (c) 2025 Bob Stammers
 *
 *
 * This file is part of IBAUK-SCOREMASTER.
 *
 * IBAUK-SCOREMASTER is free software: you can redistribute it and/or modify
 * it under the terms of the MIT License
 *
 * IBAUK-SCOREMASTER is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * MIT License for more details.
 *
 *
 * The output structure is:-
 *
 * targetFolder
 *	sm
 *		ebcimg
 *		help
 *		images
 *			bonuses
 *		jodit
 *		uploads
 *		vendor
 *	php
 *	caddy
 *
 */

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

var mySMFLAVOUR = "v3.4"
var srcFolder = flag.String("sm3", filepath.Join("..", "sm3"), "Path to ScoreMaster source")
var runsmFolder = flag.String("runsm", filepath.Join("..", "runsm"), "Path to RunSM folder")
var ebcfetchFolder = flag.String("ebcfetch", filepath.Join("..", "ebcfetch"), "Path to EBCFetch folder")
var smpatchFolder = flag.String("smpatch", filepath.Join("..", "smpatch"), "Path to SMPatch folder")
var phpFolder = flag.String("php", filepath.Join("C:\\", "PHP"), "Path to PHP installation (Windows only)")
var targetFolder = flag.String("target", "", "Path for new installation")
var db2Use = flag.String("db", "v", "v=virgin,r=rblr,l=live database")
var rblr = flag.Bool("rblr", false, "Include the extras for the RBLR1000")
var lang2use = flag.String("lang", "en", "Language code (en,de)")
var overwriteok = flag.Bool("ok", false, "Overwrite existing target")
var nodebug = flag.Bool("nodebug", false, "Don't produce debugsm")

var utilsFolder = "utils"
var sqlite3 string = "sqlite3"
var caddy string = "caddy"
var phpcgi string = "php-cgi" // non-windows only
var ebcfetch string = "ebcfetch"
var smpatch string = "smpatch"

var helpFolder = "help"

var mySMFILES = [...]string{
	"about.php", "admin.php", "bonuses.php",
	"certedit.php", "certificate.css", "certificate.php", "claims.php", "claimslog.php",
	"claimsphp.js", "cohorts.php", "combos.php", "common.php", "ereviews.php",
	"entrants.php", "exportxls.php", "emails.php", "fastodos.php", "fastodosphp.js",
	"favicon.ico", "importxls.php", "index.php", "legs.js", "legs.php",
	"Parsedown.php", "picklist.php", "reports.php", "readyset.php",
	"LICENSE", "reboot.css", "recalc.js", "recalc.php", "restbonuses.php", "scorex.php",
	"setup.php", "score.css", "score.js", "score.php", "scorecard.php", "scoring.php", "sm.php",
	"showhelp.php",
	"speeding.php", "teams.php", "utils.php", "timep.php", "cats.php",
	"classes.php",
}

var myLANGFILES = [...]string{
	"custom.js", "customvars.php",
}

var myIMAGES = [...]string{
	"ibauk.png", "ibauk90.png",
	"alertalert.png", "alertbike.png", "alertface.png", "alertdaylight.png",
	"alertnight.png", "alertreceipt.png", "alertrestricted.png", "alertteam.png",
}

var rblrIMAGES = [...]string{
	"ss1000.jpg", "smallpoppy.png", "rblr.png", "poppy.png",
	"rblrhead.png", "bb1500.jpg", "bbg1500.png",
	"route500AC.jpg", "route500CW.jpg",
}

const myREADME = `This folder contains a working copy of IBA ScoreMaster rally administration software.

If you wish, you can copy this folder and its subfolders into a location on your
hard drive and run the software from there.

To start the ScoreMaster server, run or double-click the file 'runsm' in this folder. You might need
to grant permission for your firewall but it will only fire up a local web browser to enable you to
access the system, from this machine or anywhere on your local network.

If you need further support please contact the author, Bob Stammers, at webmaster@ironbutt.co.uk
`

var smFolder string

func main() {

	fmt.Println()
	log.Println("MakeSM", mySMFLAVOUR, "ScoreMaster installation maker")
	flag.Parse()
	log.Println("Building for", runtime.GOOS)
	if *targetFolder == "" {
		log.Fatal("You must specify a target folder")
	}

	checkPrerequisites()

	if *overwriteok {
		zapTarget()
	}

	// Start by building the folder structure
	log.Println("Building folder structure")

	makeFolder(*targetFolder)

	smFolder = filepath.Join(*targetFolder, "sm")

	makeFolder(smFolder)
	makeFolder(filepath.Join(smFolder, "images"))
	makeFolder(filepath.Join(smFolder, "images", "bonuses"))
	makeFolder(filepath.Join(smFolder, "uploads"))
	makeFolder(filepath.Join(smFolder, "ebcimg"))
	makeFolder(filepath.Join(*targetFolder, "php"))
	makeFolder(filepath.Join(*targetFolder, "caddy"))

	// Now populate those folders
	writeReadme()
	copyPHP()
	copyDatabase()
	copySMFiles()
	copyExecs()
	copyImages()
	copyJodit()
	copyPhpPackages()
	generateHelp()
	log.Println("ScoreMaster installed in " + *targetFolder)
	fmt.Println()

}

func binexe(exename string) string {

	if runtime.GOOS == "windows" {
		return exename + ".exe"
	}
	return exename

}

func checkPrerequisites() {

	var ok = true
	var sqlitetest = binexe(filepath.Join(*srcFolder, utilsFolder, sqlite3))
	var caddytest = binexe(filepath.Join(*srcFolder, utilsFolder, caddy))
	var runtest = binexe(filepath.Join(*runsmFolder, "runsm"))
	var cgitest = binexe(filepath.Join(*srcFolder, utilsFolder, phpcgi))
	var jodittest = filepath.Join(*srcFolder, "jodit")
	var vendortest = filepath.Join(*srcFolder, "vendor")
	var ebcfetchtest = binexe(filepath.Join(*ebcfetchFolder, ebcfetch))
	var smpatchtest = binexe(filepath.Join(*smpatchFolder, smpatch))

	if runtime.GOOS == "windows" && !fileOrFolderExists(*phpFolder) {
		log.Printf("*** %s does not exist!", *phpFolder)
		log.Printf("*** You must have a working PHP installation installed. Download from php.net")
		log.Printf("*** Be sure to use my php.ini rather than the default")
		ok = false
	}

	if runtime.GOOS != "windows" && runtime.GOOS != "linux" && !fileOrFolderExists(cgitest) {
		log.Printf("*** %s does not exist!", cgitest)
		log.Printf("*** You must obtain a copy - compile PHP from source?")
		log.Printf("*** Be sure to use my php.ini rather than the default")
		ok = false
	}
	if !fileOrFolderExists(sqlitetest) {
		log.Printf("*** %s does not exist!", sqlitetest)
		log.Printf("*** Please download from sqlite.org")
		ok = false
	}
	if !fileOrFolderExists(ebcfetchtest) {
		log.Printf("*** %s does not exist!", ebcfetchtest)
		log.Printf("*** Please build or download from github")
		ok = false
	}
	if !fileOrFolderExists(smpatchtest) {
		log.Printf("*** %s does not exist!", smpatchtest)
		log.Printf("*** Please build or download from github")
		ok = false
	}
	if !fileOrFolderExists(caddytest) {
		log.Printf("*** %s does not exist!", caddytest)
		log.Printf("*** You must have a working Caddy installation. Download from github.com/caddyserver/caddy)")
		ok = false
	}
	if !fileOrFolderExists(jodittest) {
		log.Printf("*** %s does not exist!", jodittest)
		log.Printf("*** You might want to download from github.com/xdan/jodit")
		ok = false
	}
	if !fileOrFolderExists(vendortest) {
		log.Printf("*** %s does not exist!", vendortest)
		log.Printf("*** You probably need to run 'composer install' (download from getcomposer.org")
		ok = false
	}
	if !fileOrFolderExists(runtest) {
		log.Printf("*** %s does not exist!", runtest)
		log.Printf("*** You must do build it or use '-run' to point me to it")
		ok = false
	}
	if !ok {
		log.Fatal("*** Please fix these issues and try again")
	}

}

func copyDatabase() {

	log.Print("Establishing database")
	if *db2Use != "l" {

		if !loadSQL("ScoreMaster.sql") {
			log.Fatal("*** can't load ScoreMaster.sql")
		}
		if *lang2use != "en" {
			if !loadSQL("Reasons-" + *lang2use + ".sql") {
				log.Fatal("*** can't load foreign reasons")
			}
		}
		if *rblr || *db2Use == "r" {
			log.Print("Loading RBLR certificates")
			if !loadSQL("rblrcerts.sql") {
				log.Fatal("*** can't load rblrcerts.sql")
			}
		}
	} else {
		copyFile(filepath.Join(*srcFolder, "ScoreMaster.db"), filepath.Join(smFolder, "ScoreMaster.db"))
	}

	// Now make a copy to provide a simple means of starting again if necessary
	copyFile(filepath.Join(smFolder, "ScoreMaster.db"), filepath.Join(smFolder, "ScoreMaster-empty.db"))

}

func copyExec(src, dst string) {

	var xsrc = binexe(src)
	var xdst = binexe(dst)

	copyFile(xsrc, xdst)
	os.Chmod(dst, 0755)

}
func copyExecs() {

	log.Print("Copying executables")

	copyExec(filepath.Join(*srcFolder, utilsFolder, caddy), filepath.Join(*targetFolder, "caddy", caddy))
	copyExec(filepath.Join(*ebcfetchFolder, ebcfetch), filepath.Join(*targetFolder, "caddy", ebcfetch))
	copyExec(filepath.Join(*smpatchFolder, smpatch), filepath.Join(*targetFolder, smpatch))
	copyExec(filepath.Join(*runsmFolder, "runsm"), filepath.Join(*targetFolder, "runsm"))
	if !*nodebug {
		copyExec(filepath.Join(*runsmFolder, "runsm"), filepath.Join(*targetFolder, "debugsm"))
	}
}

func copyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	if err != nil {
		log.Fatal("*** can't copy file " + src)
	}
	return nBytes, err
}

func copyFolderTree(src string, dst string) error {
	var err error
	var fds []os.DirEntry
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = os.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := filepath.Join(src, fd.Name())
		dstfp := filepath.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = copyFolderTree(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if _, err = copyFile(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

func copyHelp(srcpath string, dstpath string, folder string) error {

	var err error
	var fds []os.DirEntry
	var mdre = regexp.MustCompile(`\.md`)

	src := filepath.Join(srcpath, folder)
	dst := filepath.Join(dstpath, folder)

	makeFolder(dst)

	if fds, err = os.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := filepath.Join(src, fd.Name())
		dstfp := filepath.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = copyHelp(src, dst, fd.Name()); err != nil {
				fmt.Println(err)
			}
		} else if mdre.MatchString(fd.Name()) {
			copyMarkdown(srcfp, dstfp)
		} else {
			if _, err = copyFile(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

func copyImages() {

	log.Print("Copying images")
	copyImageSet(myIMAGES[:])
	if *rblr || *db2Use == "r" {
		copyImageSet(rblrIMAGES[:])
	}
}

func copyImageSet(set []string) {

	for _, img := range set {
		_, err := copyFile(filepath.Join(*srcFolder, "images", img), filepath.Join(smFolder, "images", img))
		if err != nil {
			log.Fatalf("*** can't copy image %s (%s)", img, err)
		}
	}

}
func copyJodit() {

	log.Print("Copying Jodit WYSIWYG editor")

	makeFolder(filepath.Join(smFolder, "jodit"))

	copyFile(filepath.Join(*srcFolder, "jodit", "build", "jodit.min.js"), filepath.Join(smFolder, "jodit", "jodit.min.js"))
	copyFile(filepath.Join(*srcFolder, "jodit", "build", "jodit.min.css"), filepath.Join(smFolder, "jodit", "jodit.min.css"))
	copyFile(filepath.Join(*srcFolder, "images", "icons", "fields.png"), filepath.Join(smFolder, "jodit", "fields.png"))
	copyFile(filepath.Join(*srcFolder, "images", "icons", "borders.png"), filepath.Join(smFolder, "jodit", "borders.png"))

}

func copyMarkdown(src string, dst string) {

	dstname := strings.Replace(dst, ".md", ".hlp", -1)
	copyFile(src, dstname)

	/**
	file, err := os.Open(src)
	if err != nil {
		log.Panicf("failed reading file: %s", err)
	}
	defer file.Close()
	txt, _ := ioutil.ReadAll(file)
	html := blackfriday.MarkdownCommon(txt)
	f, err := os.Create(dstname)
	if err != nil {
		log.Panicf("failed creating file: %s", err)
	}
	w := bufio.NewWriter(f)
	w.Write(html)
	w.Flush()
	f.Close()
	**/
}

func copyPHP() {

	if runtime.GOOS == "windows" {
		log.Print("Copying PHP from " + *phpFolder)
		if err := copyFolderTree(*phpFolder, filepath.Join(*targetFolder, "php")); err != nil {
			log.Fatalf("*** FAILED copying folder: %s", err)
		}
	} else {
		cgi := filepath.Join(*srcFolder, utilsFolder, phpcgi)
		if _, err := os.Stat(cgi); err == nil {
			tgtcgi := filepath.Join(*targetFolder, "php", "php-cgi")
			copyExec(cgi, tgtcgi)

		}
	}
	ini := filepath.Join(*srcFolder, "php", "php.ini")
	if _, err := os.Stat(ini); err == nil {
		copyFile(ini, filepath.Join(*targetFolder, "php", "php.ini"))
	}

}

func copyPhpPackages() {

	log.Print("Copying PHP packages")
	if err := copyFolderTree(filepath.Join(*srcFolder, "vendor"), filepath.Join(smFolder, "vendor")); err != nil {
		log.Fatalf("*** FAILED copying folder: %s", err)
	}

}

func copySMFiles() {

	var lng string = "."

	log.Print("Copying main SM application")

	if *lang2use != "en" {
		lng = "-" + *lang2use + "."
	}
	for _, s := range mySMFILES {
		_, err := copyFile(filepath.Join(*srcFolder, s), filepath.Join(smFolder, s))
		if err != nil {
			log.Println("*** can't copy " + s)
		}
	}
	for _, s := range myLANGFILES {
		_, err := copyFile(filepath.Join(*srcFolder, strings.Replace(s, ".", lng, 1)), filepath.Join(smFolder, s))
		if err != nil {
			log.Println("*** can't copy " + s + lng)
		}
	}
}

func fileOrFolderExists(filename string) bool {
	_, err := os.Stat(filename)

	return !os.IsNotExist(err)
}

func generateHelp() {

	log.Print("Generating help")

	copyHelp(*srcFolder, filepath.Join(*targetFolder, "sm"), helpFolder)
}

func loadSQL(sqlfile string) bool {

	sql, _ := os.ReadFile(filepath.Join(*srcFolder, sqlfile))
	cmd := exec.Command(filepath.Join(*srcFolder, utilsFolder, sqlite3), filepath.Join(smFolder, "ScoreMaster.db"))
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Printf("%v\n", err)
		log.Fatal("*** can't load database from SQL")
	}
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, string(sql))
	}()
	_, err = cmd.CombinedOutput()
	return err == nil
}

func makeFolder(folder string) {

	if !establishFolder(folder) {
		log.Fatal("*** can't establish folder " + folder)
	}

}

func establishFolder(folder string) bool {

	err := os.Mkdir(folder, 0777)
	return err == nil
}

func writeReadme() {

	f, err := os.Create(filepath.Join(*targetFolder, "readme.txt"))
	if err != nil {
		return
	}

	f.WriteString(myREADME)
	f.WriteString("\n\n" + runtime.GOOS)
	f.Close()

	f, _ = os.Create(filepath.Join(*targetFolder, runtime.GOOS))
	f.Close()

}

func zapTarget() {

	log.Print("Overwriting " + *targetFolder)
	os.RemoveAll(*targetFolder)
	time.Sleep(1 * time.Second)
}
