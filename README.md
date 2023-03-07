# WadToGo

A WAD reader and extracting CLI tool using Golang.

## Usage

Invoke WadToGo from the command line referencing the executable and passing the filename of the WAD you want to read as a parameter.

```bash
./wadtogo [options] <WAD Filename>
```

#### Options

```
-musicinfo                  Print the songs names and format contained within the WAD file.
-mapsinfo                   Print the map names within the WAD file.
-lumpsinfo-dump <filename>  Dumps the WAD lumps list to the specified filename.
-musicinfo-dump <filename>  Dumps the songs names and format to the specified filename
-mapsinfo-dump  <filename>  Dumps the map names to the specified filename
```

_Example_:
```bash
$ ./wadtogo -musicinfo-dump musicinfo.txt DOOM.WAD
WAD Filename: DOOM.WAD
WAD Type: IWAD
Lumps: 2306
--------------------------
[Info] Song list dumped into musicinfo.txt
```

## Build

Create an executable binary for your machine using the following command:

```bash
go build .
```

## Run as developer

You can use the following command to quickly run the program without generating a binary.

```bash
go run . <WAD Filename>
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License

[MIT](https://choosealicense.com/licenses/mit/)