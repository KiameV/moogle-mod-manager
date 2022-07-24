# Example Mod

![Alt text](https://github.com/KiameV/moogle-mod-manager/blob/main/example/mod/preview.png?raw=true)

This example mod demonstrates how to have multiple sources for a single mod. The files that make up the mod are in the `mod` directory. `mod.json` is what'll be given to the mod manager to add the mod.

## Downloadables
Contains names and sources for 3 different sources (downloads):
Name | Source
--- | ---
base mod | base_mod.zip
normal cosmog | normal_cosmog.zip
invert cosmog | invert_cosmog.zip

## Usage
When this mod is enabled, the user will be taken to a configuration screen asking how they'd like Cosmog to look. 
The user has two options:
- Normal
- Invert

In both cases the user is given a description of their choice as well as an image of their selection.

Once the choice is made, the mod manager will download:
- `base_mod.zip` as it's specified as `Always Install`
- `normal_cosmog.zip` if the user chose Normal
- `invert_cosmog.zip` if the user chose Invert

As part of the installation the following files will be copied:
- `base_mod.zip`
  - `moogles_to_manage.txt` will be copied from `./assets/file` to `./assets`
  - `Mog.png`, `Molulu.png`, and `Mugmug.png` will be copied from directory `./assets/dir` to `./assets`
- `normal_cosmog.zip` (if chosen)
  - `Cosmog.png` from `.` to `./assets`
- `invert_cosmog.zip` (if chosen)
  - `Cosmog.png` from `.` to `./assets`

## Considerations
If I decided I only wanted a single downloadable I can specify in each configuration choice to use 
`base mod` instead. From there I'd change the file mapping to copy from a different place.
