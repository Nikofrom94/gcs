/*
 * Copyright ©1998-2022 by Richard A. Wilkes. All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, version 2.0. If a copy of the MPL was not distributed with
 * this file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 * This Source Code Form is "Incompatible With Secondary Licenses", as
 * defined by the Mozilla Public License, version 2.0.
 */

package ux

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/richardwilkes/gcs/v5/model"
	"github.com/richardwilkes/gcs/v5/svg"
	"github.com/richardwilkes/toolbox"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/toolbox/xio/fs"
	"github.com/richardwilkes/unison"
)

/*
var (
	_        FileBackedDockable         = &VehiculeSheet{}
	_        unison.UndoManagerProvider = &VehiculeSheet{}
	_        ModifiableRoot             = &VehiculeSheet{}
	_        Rebuildable                = &VehiculeSheet{}
	_        unison.TabCloser           = &VehiculeSheet{}
	dropKeys                            = []string{
		equipmentDragKey,
		model.SkillID,
		model.SpellID,
		traitDragKey,
		noteDragKey,
	}
)
*/
// Sheet holds the view for a GURPS character sheet.
type VehiculeSheet struct {
	unison.Panel
	path                 string
	targetMgr            *TargetMgr
	undoMgr              *unison.UndoManager
	toolbar              *unison.Panel
	scroll               *unison.ScrollPanel
	vehicule             *model.Vehicule
	crc                  uint64
	content              *unison.Panel
	modifiedFunc         func()
	Reactions            *PageList[*model.ConditionalModifier]
	ConditionalModifiers *PageList[*model.ConditionalModifier]
	MeleeWeapons         *PageList[*model.Weapon]
	RangedWeapons        *PageList[*model.Weapon]
	Traits               *PageList[*model.Trait]
	Skills               *PageList[*model.Skill]
	Spells               *PageList[*model.Spell]
	CarriedEquipment     *PageList[*model.Equipment]
	OtherEquipment       *PageList[*model.Equipment]
	Notes                *PageList[*model.Note]
	dragReroutePanel     *unison.Panel
	scale                int
	awaitingUpdate       bool
	needsSaveAsPrompt    bool
}

// ActiveVehiculeSheet returns the currently active sheet.
func ActiveVehiculeSheet() *VehiculeSheet {
	d := ActiveDockable()
	if d == nil {
		return nil
	}
	if s, ok := d.(*VehiculeSheet); ok {
		return s
	}
	return nil
}

// OpenVehiculeSheets returns the currently open sheets.
func OpenVehiculeSheets(exclude *VehiculeSheet) []*VehiculeSheet {
	var sheets []*VehiculeSheet
	ws := WorkspaceFromWindowOrAny(unison.ActiveWindow())
	ws.DocumentDock.RootDockLayout().ForEachDockContainer(func(dc *unison.DockContainer) bool {
		for _, one := range dc.Dockables() {
			if sheet, ok := one.(*VehiculeSheet); ok && sheet != exclude {
				sheets = append(sheets, sheet)
			}
		}
		return false
	})
	return sheets
}

// NewVehiculeSheetSheetFromFile loads a GURPS vehicule sheet file and creates a new unison.Dockable for it.
func NewVehiculeSheetSheetFromFile(filePath string) (unison.Dockable, error) {
	vehicule, err := model.NewVehiculeFromFile(os.DirFS(filepath.Dir(filePath)), filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	s := NewVehiculeSheet(filePath, vehicule)
	s.needsSaveAsPrompt = false
	return s, nil
}

func (vs *VehiculeSheet) createVehiculePageTopBlock(targetMgr *TargetMgr) (page *VehiculePage, modifiedFunc func()) {
	page = NewVehiculePage(vs.vehicule)
	var top *unison.Panel
	top, modifiedFunc = createVehiculePageFirstRow(vs.vehicule, targetMgr)
	page.AddChild(top)
	page.AddChild(createVehiculePageSecondRow(vs.vehicule, targetMgr))
	return page, modifiedFunc
}

// NewVehiculeSheet creates a new unison.Dockable for GURPS character sheet files.
func NewVehiculeSheet(filePath string, vehicule *model.Vehicule) *VehiculeSheet {
	s := &VehiculeSheet{
		path:              filePath,
		undoMgr:           unison.NewUndoManager(200, func(err error) { jot.Error(err) }),
		scroll:            unison.NewScrollPanel(),
		vehicule:          vehicule,
		crc:               vehicule.CRC64(),
		scale:             model.GlobalSettings().General.InitialSheetUIScale,
		content:           unison.NewPanel(),
		needsSaveAsPrompt: true,
	}
	s.Self = s
	s.targetMgr = NewTargetMgr(s)
	s.SetLayout(&unison.FlexLayout{
		Columns: 1,
		HAlign:  unison.FillAlignment,
		VAlign:  unison.FillAlignment,
	})

	s.MouseDownCallback = func(_ unison.Point, _, _ int, _ unison.Modifiers) bool {
		s.RequestFocus()
		return false
	}
	s.DataDragOverCallback = func(_ unison.Point, data map[string]any) bool {
		s.dragReroutePanel = nil
		for _, key := range dropKeys {
			if _, ok := data[key]; ok {
				if s.dragReroutePanel = s.keyToPanel(key); s.dragReroutePanel != nil {
					s.dragReroutePanel.DataDragOverCallback(unison.Point{Y: 100000000}, data)
					return true
				}
				break
			}
		}
		return false
	}
	s.DataDragExitCallback = func() {
		if s.dragReroutePanel != nil {
			s.dragReroutePanel.DataDragExitCallback()
			s.dragReroutePanel = nil
		}
	}
	s.DataDragDropCallback = func(_ unison.Point, data map[string]any) {
		if s.dragReroutePanel != nil {
			s.dragReroutePanel.DataDragDropCallback(unison.Point{Y: 10000000}, data)
			s.dragReroutePanel = nil
		}
	}
	s.DrawOverCallback = func(gc *unison.Canvas, rect unison.Rect) {
		if s.dragReroutePanel != nil {
			r := s.RectFromRoot(s.dragReroutePanel.RectToRoot(s.dragReroutePanel.ContentRect(true)))
			paint := unison.DropAreaColor.Paint(gc, r, unison.Fill)
			paint.SetColorFilter(unison.Alpha30Filter())
			gc.DrawRect(r, paint)
		}
	}

	s.content.SetLayout(&unison.FlexLayout{
		Columns:  1,
		VSpacing: 1,
	})
	var vPage *VehiculePage
	vPage, s.modifiedFunc = createVehiculePageTopBlock(s.vehicule, s.targetMgr)
	s.content.AddChild(vPage)
	s.createLists()
	s.scroll.SetContent(s.content, unison.UnmodifiedBehavior, unison.UnmodifiedBehavior)
	s.scroll.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		VAlign: unison.FillAlignment,
		HGrab:  true,
		VGrab:  true,
	})
	s.scroll.DrawCallback = func(gc *unison.Canvas, rect unison.Rect) {
		gc.DrawRect(rect, model.PageVoidColor.Paint(gc, rect, unison.Fill))
	}

	helpButton := unison.NewSVGButton(svg.Help)
	helpButton.Tooltip = unison.NewTooltipWithText(i18n.Text("Help"))
	helpButton.ClickCallback = func() { HandleLink(nil, "md:Help/Interface/Character Sheet") }

	sheetSettingsButton := unison.NewSVGButton(svg.Settings)
	sheetSettingsButton.Tooltip = unison.NewTooltipWithText(i18n.Text("Sheet Settings"))
	//sheetSettingsButton.ClickCallback = func() { ShowSheetSettings(s) }

	attributesButton := unison.NewSVGButton(svg.Attributes)
	attributesButton.Tooltip = unison.NewTooltipWithText(i18n.Text("Attributes"))
	//attributesButton.ClickCallback = func() { ShowAttributeSettings(s) }

	/* no calculator for vehicule sheet
	calcButton := unison.NewSVGButton(svg.Calculator)
	calcButton.Tooltip = unison.NewTooltipWithText(i18n.Text("Calculators (jumping, throwing, hiking, etc.)"))
	calcButton.ClickCallback = func() { DisplayCalculator(s) }
	*/
	s.toolbar = unison.NewPanel()
	s.toolbar.SetBorder(unison.NewCompoundBorder(unison.NewLineBorder(unison.DividerColor, 0, unison.Insets{Bottom: 1},
		false), unison.NewEmptyBorder(unison.StdInsets())))
	s.toolbar.SetLayoutData(&unison.FlexLayoutData{
		HAlign: unison.FillAlignment,
		HGrab:  true,
	})
	s.toolbar.AddChild(NewDefaultInfoPop())
	s.toolbar.AddChild(helpButton)
	s.toolbar.AddChild(
		NewScaleField(
			model.InitialUIScaleMin,
			model.InitialUIScaleMax,
			func() int { return model.GlobalSettings().General.InitialSheetUIScale },
			func() int { return s.scale },
			func(scale int) { s.scale = scale },
			nil,
			false,
			s.scroll,
		),
	)
	s.toolbar.AddChild(sheetSettingsButton)
	s.toolbar.AddChild(attributesButton)
	//s.toolbar.AddChild(bodyTypeButton)
	s.toolbar.AddChild(NewToolbarSeparator())
	//s.toolbar.AddChild(calcButton)
	s.toolbar.AddChild(NewToolbarSeparator())
	installSearchTracker(s.toolbar, func() {
		s.Reactions.Table.ClearSelection()
		s.ConditionalModifiers.Table.ClearSelection()
		s.MeleeWeapons.Table.ClearSelection()
		s.RangedWeapons.Table.ClearSelection()
		s.Traits.Table.ClearSelection()
		s.Skills.Table.ClearSelection()
		s.Spells.Table.ClearSelection()
		s.CarriedEquipment.Table.ClearSelection()
		s.OtherEquipment.Table.ClearSelection()
		s.Notes.Table.ClearSelection()
	}, func(refList *[]*searchRef, text string) {
		searchSheetTable(refList, text, s.Traits)
		searchSheetTable(refList, text, s.Skills)
		searchSheetTable(refList, text, s.Spells)
		searchSheetTable(refList, text, s.CarriedEquipment)
		searchSheetTable(refList, text, s.OtherEquipment)
		searchSheetTable(refList, text, s.Notes)
	})
	s.toolbar.SetLayout(&unison.FlexLayout{
		Columns:  len(s.toolbar.Children()),
		HSpacing: unison.StdHSpacing,
	})

	s.AddChild(s.toolbar)
	s.AddChild(s.scroll)

	s.InstallCmdHandlers(SaveItemID, func(_ any) bool { return s.Modified() }, func(_ any) { s.save(false) })
	s.InstallCmdHandlers(SaveAsItemID, unison.AlwaysEnabled, func(_ any) { s.save(true) })
	s.installNewItemCmdHandlers(NewNoteItemID, NewNoteContainerItemID, s.Notes)
	s.InstallCmdHandlers(SwapDefaultsItemID, s.canSwapDefaults, s.swapDefaults)
	s.InstallCmdHandlers(ExportAsPDFItemID, unison.AlwaysEnabled, func(_ any) { s.exportToPDF() })
	s.InstallCmdHandlers(ExportAsWEBPItemID, unison.AlwaysEnabled, func(_ any) { s.exportToWEBP() })
	s.InstallCmdHandlers(ExportAsPNGItemID, unison.AlwaysEnabled, func(_ any) { s.exportToPNG() })
	s.InstallCmdHandlers(ExportAsJPEGItemID, unison.AlwaysEnabled, func(_ any) { s.exportToJPEG() })
	s.InstallCmdHandlers(PrintItemID, unison.AlwaysEnabled, func(_ any) { s.print() })
	s.InstallCmdHandlers(ClearPortraitItemID, s.canClearPortrait, s.clearPortrait)

	return s
}

func (s *VehiculeSheet) canClearPortrait(_ any) bool {
	return len(s.vehicule.Profile.PortraitData) != 0
}

func (s *VehiculeSheet) clearPortrait(_ any) {
	if s.canClearPortrait(nil) {
		s.undoMgr.Add(&unison.UndoEdit[[]byte]{
			ID:         unison.NextUndoID(),
			EditName:   clearPortraitAction.Title,
			UndoFunc:   func(edit *unison.UndoEdit[[]byte]) { s.updatePortrait(edit.BeforeData) },
			RedoFunc:   func(edit *unison.UndoEdit[[]byte]) { s.updatePortrait(edit.AfterData) },
			BeforeData: s.vehicule.Profile.PortraitData,
			AfterData:  nil,
		})
		s.updatePortrait(nil)
	}
}

func (s *VehiculeSheet) updatePortrait(data []byte) {
	s.vehicule.Profile.PortraitData = data
	s.vehicule.Profile.PortraitImage = nil
	s.MarkForRedraw()
	s.MarkModified(s)
}

func (s *VehiculeSheet) keyToPanel(key string) *unison.Panel {
	var p unison.Paneler
	switch key {
	case equipmentDragKey:
		p = s.CarriedEquipment.Table
	case model.SkillID:
		p = s.Skills.Table
	case model.SpellID:
		p = s.Spells.Table
	case traitDragKey:
		p = s.Traits.Table
	case noteDragKey:
		p = s.Notes.Table
	default:
		return nil
	}
	return p.AsPanel()
}

func (s *VehiculeSheet) installNewItemCmdHandlers(itemID, containerID int, creator itemCreator) {
	variant := NoItemVariant
	if containerID == -1 {
		variant = AlternateItemVariant
	} else {
		s.InstallCmdHandlers(containerID, unison.AlwaysEnabled,
			func(_ any) { creator.CreateItem(s, ContainerItemVariant) })
	}
	s.InstallCmdHandlers(itemID, unison.AlwaysEnabled, func(_ any) { creator.CreateItem(s, variant) })
}

// DockableKind implements widget.DockableKind
func (s *VehiculeSheet) DockableKind() string {
	return SheetDockableKind
}

// Vehicule returns the entity this is displaying information for.
func (s *VehiculeSheet) Entity() *model.Vehicule {
	return s.vehicule
}

// UndoManager implements undo.Provider
func (s *VehiculeSheet) UndoManager() *unison.UndoManager {
	return s.undoMgr
}

// TitleIcon implements workspace.FileBackedDockable
func (s *VehiculeSheet) TitleIcon(suggestedSize unison.Size) unison.Drawable {
	return &unison.DrawableSVG{
		SVG:  model.FileInfoFor(s.path).SVG,
		Size: suggestedSize,
	}
}

// Title implements workspace.FileBackedDockable
func (s *VehiculeSheet) Title() string {
	return fs.BaseName(s.path)
}

func (s *VehiculeSheet) String() string {
	return s.Title()
}

// Tooltip implements workspace.FileBackedDockable
func (s *VehiculeSheet) Tooltip() string {
	return s.path
}

// BackingFilePath implements workspace.FileBackedDockable
func (s *VehiculeSheet) BackingFilePath() string {
	return s.path
}

// SetBackingFilePath implements workspace.FileBackedDockable
func (s *VehiculeSheet) SetBackingFilePath(p string) {
	s.path = p
	if dc := unison.Ancestor[*unison.DockContainer](s); dc != nil {
		dc.UpdateTitle(s)
	}
}

// Modified implements workspace.FileBackedDockable
func (s *VehiculeSheet) Modified() bool {
	return s.crc != s.vehicule.CRC64()
}

// MarkModified implements widget.ModifiableRoot.
func (s *VehiculeSheet) MarkModified(src unison.Paneler) {
	if !s.awaitingUpdate {
		s.awaitingUpdate = true
		h, v := s.scroll.Position()
		focusRefKey := s.targetMgr.CurrentFocusRef()
		s.vehicule.DiscardCaches()
		s.modifiedFunc()
		// TODO: This can be too slow when the lists have many rows of content, impinging upon interactive typing.
		//       Looks like most of the time is spent in updating the tables. Unfortunately, there isn't a fast way to
		//       determine that the content of a table doesn't need to be refreshed.
		skipDeepSync := false
		if !toolbox.IsNil(src) {
			_, skipDeepSync = src.AsPanel().ClientData()[SkipDeepSync]
		}
		if !skipDeepSync {
			DeepSync(s)
		}
		if dc := unison.Ancestor[*unison.DockContainer](s); dc != nil {
			dc.UpdateTitle(s)
		}
		s.awaitingUpdate = false
		s.targetMgr.ReacquireFocus(focusRefKey, s.toolbar, s.scroll.Content())
		s.scroll.SetPosition(h, v)
		//UpdateCalculator(s)
	}
}

// MayAttemptClose implements unison.TabCloser
func (s *VehiculeSheet) MayAttemptClose() bool {
	return MayAttemptCloseOfGroup(s)
}

// AttemptClose implements unison.TabCloser
func (s *VehiculeSheet) AttemptClose() bool {
	if !CloseGroup(s) {
		return false
	}
	if s.Modified() {
		switch unison.YesNoCancelDialog(fmt.Sprintf(i18n.Text("Save changes made to\n%s?"), s.Title()), "") {
		case unison.ModalResponseDiscard:
		case unison.ModalResponseOK:
			if !s.save(false) {
				return false
			}
		case unison.ModalResponseCancel:
			return false
		}
	}
	if dc := unison.Ancestor[*unison.DockContainer](s); dc != nil {
		dc.Close(s)
	}
	return true
}

func (s *VehiculeSheet) save(forceSaveAs bool) bool {
	success := false
	if forceSaveAs || s.needsSaveAsPrompt {
		success = SaveDockableAs(s, model.SheetExt, s.vehicule.Save, func(path string) {
			s.crc = s.vehicule.CRC64()
			s.path = path
		})
	} else {
		success = SaveDockable(s, s.vehicule.Save, func() { s.crc = s.vehicule.CRC64() })
	}
	if success {
		s.needsSaveAsPrompt = false
	}
	return success
}

func (s *VehiculeSheet) print() {
	/*
		data, err := newPageExporter(s.entity).exportAsPDFBytes()
		if err != nil {
			unison.ErrorDialogWithError(i18n.Text("Unable to create PDF!"), err)
			return
		}
		dialog := PrintMgr.NewJobDialog(printing.PrinterID{}, "application/pdf", nil)
		if dialog.RunModal() {
			go backgroundPrint(s.vehicule.Profile.Name, dialog.Printer(), dialog.JobAttributes(), data)
		}
	*/
}

func (s *VehiculeSheet) exportToPDF() {
	/*
		s.Window().ShowCursor()
		dialog := unison.NewSaveDialog()
		dialog.SetInitialDirectory(filepath.Dir(s.BackingFilePath()))
		dialog.SetAllowedExtensions("pdf")
		if dialog.RunModal() {
			if filePath, ok := unison.ValidateSaveFilePath(dialog.Path(), "pdf", false); ok {
				model.GlobalSettings().SetLastDir(model.DefaultLastDirKey, filepath.Dir(filePath))
				if err := newPageExporter(s.entity).exportAsPDFFile(filePath); err != nil {
					unison.ErrorDialogWithError(i18n.Text("Unable to export as PDF!"), err)
				}
			}
		}
	*/
}

func (s *VehiculeSheet) exportToWEBP() {
	/*
		s.Window().ShowCursor()
		dialog := unison.NewSaveDialog()
		dialog.SetInitialDirectory(filepath.Dir(s.BackingFilePath()))
		dialog.SetAllowedExtensions("webp")
		if dialog.RunModal() {
			if filePath, ok := unison.ValidateSaveFilePath(dialog.Path(), "webp", false); ok {
				model.GlobalSettings().SetLastDir(model.DefaultLastDirKey, filepath.Dir(filePath))
				if err := newPageExporter(s.entity).exportAsWEBPs(filePath); err != nil {
					unison.ErrorDialogWithError(i18n.Text("Unable to export as WEBP!"), err)
				}
			}
		}
	*/
}

func (s *VehiculeSheet) exportToPNG() {
	/*
		s.Window().ShowCursor()
		dialog := unison.NewSaveDialog()
		dialog.SetInitialDirectory(filepath.Dir(s.BackingFilePath()))
		dialog.SetAllowedExtensions("png")
		if dialog.RunModal() {
			if filePath, ok := unison.ValidateSaveFilePath(dialog.Path(), "png", false); ok {
				model.GlobalSettings().SetLastDir(model.DefaultLastDirKey, filepath.Dir(filePath))
				if err := newPageExporter(s.entity).exportAsPNGs(filePath); err != nil {
					unison.ErrorDialogWithError(i18n.Text("Unable to export as PNG!"), err)
				}
			}
		}
	*/
}

func (s *VehiculeSheet) exportToJPEG() {
	/*
		s.Window().ShowCursor()
		dialog := unison.NewSaveDialog()
		dialog.SetInitialDirectory(filepath.Dir(s.BackingFilePath()))
		dialog.SetAllowedExtensions("jpeg")
		if dialog.RunModal() {
			if filePath, ok := unison.ValidateSaveFilePath(dialog.Path(), "jpeg", false); ok {
				model.GlobalSettings().SetLastDir(model.DefaultLastDirKey, filepath.Dir(filePath))
				if err := newPageExporter(s.entity).exportAsJPEGs(filePath); err != nil {
					unison.ErrorDialogWithError(i18n.Text("Unable to export as JPEG!"), err)
				}
			}
		}
	*/
}

func (s *VehiculeSheet) createLists() {
	children := s.content.Children()
	if len(children) == 0 {
		return
	}
	page, ok := children[0].Self.(*Page)
	if !ok {
		return
	}
	children = page.Children()
	if len(children) < 2 {
		return
	}
	for i := len(children) - 1; i > 1; i-- {
		page.RemoveChildAtIndex(i)
	}
	// Add the various blocks, based on the layout preference.
	for _, col := range s.vehicule.SheetSettings.BlockLayout.ByRow() {
		rowPanel := unison.NewPanel()
		rowPanel.SetLayout(&unison.FlexLayout{
			Columns:      len(col),
			HSpacing:     1,
			HAlign:       unison.FillAlignment,
			VAlign:       unison.FillAlignment,
			EqualColumns: true,
		})
		rowPanel.SetLayoutData(&unison.FlexLayoutData{
			HAlign: unison.FillAlignment,
			HGrab:  true,
		})
		for _, c := range col {
			switch c {
			case model.BlockLayoutSpaceshipHullsKey:
				if s.Reactions == nil {
					s.Reactions = NewReactionsPageList(s.entity)
				} else {
					s.Reactions.Sync()
				}
				rowPanel.AddChild(s.Reactions)

			case model.BlockLayoutNotesKey:
				if s.Notes == nil {
					s.Notes = NewNotesPageList(s, s.entity)
				} else {
					s.Notes.Sync()
				}
				rowPanel.AddChild(s.Notes)
			}
		}

		page.AddChild(rowPanel)
	}
	page.ApplyPreferredSize()
}

func (s *VehiculeSheet) canSwapDefaults(_ any) bool {
	canSwap := false
	for _, skillNode := range s.Skills.SelectedNodes(true) {
		skill := skillNode.Data()
		if skill.Type == model.TechniqueID {
			return false
		}
		if !skill.CanSwapDefaultsWith(skill.DefaultSkill()) && skill.BestSwappableSkill() == nil {
			return false
		}
		canSwap = true
	}
	return canSwap
}

func (s *VehiculeSheet) swapDefaults(_ any) {
	undo := &unison.UndoEdit[*TableUndoEditData[*model.Skill]]{
		ID:       unison.NextUndoID(),
		EditName: i18n.Text("Swap Defaults"),
		UndoFunc: func(e *unison.UndoEdit[*TableUndoEditData[*model.Skill]]) { e.BeforeData.Apply() },
		RedoFunc: func(e *unison.UndoEdit[*TableUndoEditData[*model.Skill]]) { e.AfterData.Apply() },
		AbsorbFunc: func(e *unison.UndoEdit[*TableUndoEditData[*model.Skill]], other unison.Undoable) bool {
			return false
		},
		BeforeData: NewTableUndoEditData(s.Skills.Table),
	}
	for _, skillNode := range s.Skills.SelectedNodes(true) {
		skill := skillNode.Data()
		if !skill.CanSwapDefaults() {
			continue
		}
		swap := skill.DefaultSkill()
		if !skill.CanSwapDefaultsWith(swap) {
			swap = skill.BestSwappableSkill()
		}
		skill.DefaultedFrom = nil
		swap.SwapDefaults()
	}
	s.Skills.Sync()
	undo.AfterData = NewTableUndoEditData(s.Skills.Table)
	s.UndoManager().Add(undo)
}

// SheetSettingsUpdated implements gurps.SheetSettingsResponder.
func (s *VehiculeSheet) SheetSettingsUpdated(vehicule *model.Vehicule, blockLayout bool) {
	if s.vehicule == vehicule {
		s.MarkModified(nil)
		s.Rebuild(blockLayout)
	}
}

// Rebuild implements widget.Rebuildable.
func (s *VehiculeSheet) Rebuild(full bool) {
	h, v := s.scroll.Position()
	focusRefKey := s.targetMgr.CurrentFocusRef()

	if full {

		notesSelMap := s.Notes.RecordSelection()
		defer func() {

			s.Notes.ApplySelection(notesSelMap)
		}()
		s.createLists()
	}
	DeepSync(s)
	if dc := unison.Ancestor[*unison.DockContainer](s); dc != nil {
		dc.UpdateTitle(s)
	}
	s.targetMgr.ReacquireFocus(focusRefKey, s.toolbar, s.scroll.Content())
	s.scroll.SetPosition(h, v)

}
