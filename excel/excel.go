package excel

import (
	"github.com/xuri/excelize/v2"
	"github.com/yungsem/db-desc/db"
	"strconv"
)

func GenerateExcel(tableInfos []db.TableInfo) error {
	f := excelize.NewFile()

	sheet := "表结构说明"
	excelName := "MOM表结构说明"
	_, err := f.NewSheet(sheet)
	if err != nil {
		return err
	}

	// 创建目录的 sheet
	indexSheet := "目录"
	// 将默认的名为 Sheet1 的工作表重名为配置文件中指定的名称
	err = f.SetSheetName("Sheet1", indexSheet)
	if err != nil {
		return err
	}

	styleWithColorBoldId, err := f.NewStyle(styleWithColor(true, true))
	if err != nil {
		return err
	}

	styleId, err := f.NewStyle(style(false, true))
	if err != nil {
		return err
	}

	rowIndex := 0
	for i, ti := range tableInfos {
		// 写入行：表名
		rowIndex++
		i++
		err = f.SetSheetRow(sheet, "A"+strconv.Itoa(rowIndex), &[]interface{}{"表名", ti.TableName})
		if err != nil {
			return err
		}

		err = f.SetSheetRow(indexSheet, "A"+strconv.Itoa(i), &[]interface{}{ti.TableName, ti.TableComment.String})
		if err != nil {
			return err
		}

		err = f.SetCellHyperLink(indexSheet, "A"+strconv.Itoa(i), sheet+"!"+"A"+strconv.Itoa(rowIndex), "Location")
		if err != nil {
			return err
		}

		indexStyle, err := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{Color: "1265BE", Underline: "single"},
		})
		if err != nil {
			return err
		}

		err = f.SetCellStyle(indexSheet, "A"+strconv.Itoa(i), "A"+strconv.Itoa(i), indexStyle)
		if err != nil {
			return err
		}
		err = f.SetColWidth(indexSheet, "A", "B", 80)
		if err != nil {
			return err
		}

		// 设置行样式
		err = f.SetCellStyle(sheet, "A"+strconv.Itoa(rowIndex), "H"+strconv.Itoa(rowIndex), styleWithColorBoldId)
		if err != nil {
			return err
		}
		// 合并单元格
		err = f.MergeCell(sheet, "B"+strconv.Itoa(rowIndex), "H"+strconv.Itoa(rowIndex))
		if err != nil {
			return err
		}
		// 写入行：表备注
		rowIndex++
		err = f.SetSheetRow(sheet, "A"+strconv.Itoa(rowIndex), &[]interface{}{"表备注", ti.TableComment.String})
		if err != nil {
			return err
		}
		// 设置行样式
		err = f.SetCellStyle(sheet, "A"+strconv.Itoa(rowIndex), "H"+strconv.Itoa(rowIndex), styleWithColorBoldId)
		if err != nil {
			return err
		}
		// 合并单元格
		err = f.MergeCell(sheet, "B"+strconv.Itoa(rowIndex), "H"+strconv.Itoa(rowIndex))
		if err != nil {
			return err
		}

		// 列内容开始
		// 设置行样式
		err = f.SetCellStyle(sheet, "A"+strconv.Itoa(rowIndex+1), "H"+strconv.Itoa(rowIndex), styleWithColorBoldId)
		if err != nil {
			return err
		}

		// 写入行：列头
		rowIndex++
		err = f.SetSheetRow(sheet, "A"+strconv.Itoa(rowIndex),
			&[]interface{}{"列名", "列类型", "长度", "精度", "允许为空", "默认值", "说明", "是否是主键"})
		if err != nil {
			return err
		}
		for _, ci := range ti.ColumnInfos {
			rowIndex++
			arr := []interface{}{
				ci.Name,
				ci.Kind,
				ci.Length.String,
				ci.Precision.String,
				ci.NullFlag,
				ci.DefaultValue.String,
				ci.Comment.String,
				ci.PkFlag,
			}
			err = f.SetSheetRow(sheet, "A"+strconv.Itoa(rowIndex), &arr)
			if err != nil {
				return err
			}
			// 设置行样式
			err = f.SetCellStyle(sheet, "A"+strconv.Itoa(rowIndex), "H"+strconv.Itoa(rowIndex), styleId)
			if err != nil {
				return err
			}
		}

		// 每个表结束空一行
		rowIndex++
	}

	// 统一设置列宽
	err = f.SetColWidth(sheet, "A", "A", 30)
	if err != nil {
		return err
	}
	err = f.SetColWidth(sheet, "B", "B", 15)
	if err != nil {
		return err
	}
	err = f.SetColWidth(sheet, "C", "C", 10)
	if err != nil {
		return err
	}
	err = f.SetColWidth(sheet, "D", "D", 10)
	if err != nil {
		return err
	}
	err = f.SetColWidth(sheet, "E", "E", 10)
	if err != nil {
		return err
	}
	err = f.SetColWidth(sheet, "F", "F", 10)
	if err != nil {
		return err
	}
	err = f.SetColWidth(sheet, "G", "G", 50)
	if err != nil {
		return err
	}
	err = f.SetColWidth(sheet, "H", "H", 10)
	if err != nil {
		return err
	}

	// 保存为 excel 文件
	err = f.SaveAs(excelName + ".xlsx")
	if err != nil {
		return err
	}

	return nil
}

func setIndex(f *excelize.File, indexSheet string, tableSheet string, indexList []string) error {
	for i, v := range indexList {
		i++
		err := f.SetSheetRow(indexSheet, "A"+strconv.Itoa(i), &[]interface{}{v})
		if err != nil {
			return err
		}

		err = f.SetCellHyperLink(indexSheet, "A"+strconv.Itoa(i), tableSheet+"!"+v, "Location")
		if err != nil {
			return err
		}
	}
	return nil
}

// styleWithColor 创建带底色的单元格样式
func styleWithColor(bold bool, autoWrap bool) *excelize.Style {
	s := style(bold, autoWrap)
	s.Fill = excelize.Fill{
		Type:    "pattern",
		Pattern: 1,
		Color:   []string{"#BDD7EE"},
		Shading: 0,
	}
	return s
}

// style 创建不带底色的单元格样式
func style(bold bool, autoWrap bool) *excelize.Style {
	return &excelize.Style{
		Border: []excelize.Border{
			{
				Type:  "left",
				Color: "000000",
				Style: 1,
			},
			{
				Type:  "right",
				Color: "000000",
				Style: 1,
			},
			{
				Type:  "top",
				Color: "000000",
				Style: 1,
			},
			{
				Type:  "bottom",
				Color: "000000",
				Style: 1,
			},
		},
		Font: &excelize.Font{
			Bold:         bold,
			Italic:       false,
			Underline:    "",
			Family:       "Microsoft YaHei",
			Size:         10,
			Strike:       false,
			Color:        "000000",
			ColorIndexed: 0,
			ColorTheme:   nil,
			ColorTint:    0,
			VertAlign:    "",
		},
		Alignment: &excelize.Alignment{
			WrapText: autoWrap,
			Vertical: "center",
		},
	}
}
