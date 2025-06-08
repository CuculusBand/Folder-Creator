# Folder-Creator
A tool that reads folder names from a CSV or Excel file and creates folders automatically.

For each row:
1. The first column of each row defines a top-level folder.
2. The remaining columns (if any) define subfolders under that top-level folder.
3. All subfolders are created at the same level, not nested within each other.

---------------------------------------

Example:

Input table

| A        | B           | C      |
| -------- | ----------- | ------ |
| Project1 | Data        | Report |
| Project2 | Data        |        |
| Project3 | Note        | Review |

Output structure
```
Project1/
├── Data/
└── Reprot/
Project2/
└── Data/
Project3/
├── Note/
└── Review/
```