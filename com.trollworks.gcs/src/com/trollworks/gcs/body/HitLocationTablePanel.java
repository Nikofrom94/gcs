/*
 * Copyright ©1998-2021 by Richard A. Wilkes. All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, version 2.0. If a copy of the MPL was not distributed with
 * this file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 * This Source Code Form is "Incompatible With Secondary Licenses", as
 * defined by the Mozilla Public License, version 2.0.
 */

package com.trollworks.gcs.body;

import com.trollworks.gcs.character.FieldFactory;
import com.trollworks.gcs.ui.Colors;
import com.trollworks.gcs.ui.DynamicColor;
import com.trollworks.gcs.ui.ThemeColor;
import com.trollworks.gcs.ui.border.LineBorder;
import com.trollworks.gcs.ui.layout.PrecisionLayout;
import com.trollworks.gcs.ui.layout.PrecisionLayoutData;
import com.trollworks.gcs.ui.widget.BandedPanel;
import com.trollworks.gcs.ui.widget.EditorField;
import com.trollworks.gcs.ui.widget.FontAwesomeButton;
import com.trollworks.gcs.ui.widget.Label;
import com.trollworks.gcs.ui.widget.Wrapper;
import com.trollworks.gcs.utility.Dice;
import com.trollworks.gcs.utility.I18n;
import com.trollworks.gcs.utility.ID;
import com.trollworks.gcs.utility.text.DiceFormatter;
import com.trollworks.gcs.utility.text.Text;

import java.awt.Color;
import java.awt.Component;
import java.awt.Container;
import java.awt.EventQueue;
import java.awt.Rectangle;
import javax.swing.JFormattedTextField;
import javax.swing.SwingConstants;
import javax.swing.text.DefaultFormatterFactory;

public class HitLocationTablePanel extends BandedPanel {
    private HitLocationTable mLocations;
    private Runnable         mAdjustCallback;
    private EditorField      mFirstField;

    public HitLocationTablePanel(HitLocationTable locations, Runnable adjustCallback) {
        super(false);
        setLayout(new PrecisionLayout().setMargins(0, 10, 0, 10));
        mLocations = locations;
        mAdjustCallback = adjustCallback;
        if (isSubTable()) {
            setBorder(new LineBorder(ThemeColor.DIVIDER));
            setBackground(new DynamicColor(() -> Colors.adjustSaturation(ThemeColor.BANDING, -0.05f * countSubTableDepth()).getRGB()));
        }
        fill();
    }

    public int countSubTableDepth() {
        int         depth = 0;
        HitLocation loc   = mLocations.getOwningLocation();
        while (loc != null) {
            depth++;
            HitLocationTable table = loc.getOwningTable();
            loc = table != null ? table.getOwningLocation() : null;
        }
        return depth;
    }

    public HitLocationTable getHitLocations() {
        return mLocations;
    }

    public Runnable getAdjustCallback() {
        return mAdjustCallback;
    }

    @Override
    public int getScrollableUnitIncrement(Rectangle visibleRect, int orientation, int direction) {
        return 16;
    }

    @Override
    public boolean getScrollableTracksViewportWidth() {
        return true;
    }

    public void reset(HitLocationTable locations) {
        removeAll();
        mLocations.resetTo(locations);
        fill();
    }

    private void fill() {
        mFirstField = null;
        Wrapper wrapper = new Wrapper(new PrecisionLayout().setColumns(isSubTable() ? 5 : 7).setMargins(0));
        wrapper.add(new FontAwesomeButton("\uf055", I18n.text("Add Hit Location"), this::addHitLocation));
        if (isSubTable()) {
            wrapper.add(new Label(I18n.text("Sub-Table")), new PrecisionLayoutData().setFillHorizontalAlignment());
        } else {
            mFirstField = addField(wrapper,
                    I18n.text("ID"),
                    I18n.text("An ID for the hit location table"),
                    mLocations.getID(),
                    Text.makeFiller(8, 'm'),
                    FieldFactory.STRING,
                    (f) -> {
                        String existingID = mLocations.getID();
                        String id         = ((String) f.getValue());
                        if (!existingID.equals(id)) {
                            id = ID.sanitize(id, null, false);
                            if (id.isEmpty()) {
                                f.setValue(existingID);
                            } else {
                                mLocations.setID(id);
                                f.setValue(id);
                                mAdjustCallback.run();
                            }
                        }
                    });
            addField(wrapper,
                    I18n.text("Name"),
                    I18n.text("The name of this hit location table"),
                    mLocations.getName(),
                    null,
                    FieldFactory.STRING,
                    (f) -> {
                        mLocations.setName((String) f.getValue());
                        mAdjustCallback.run();
                    });
        }
        EditorField field = addField(wrapper,
                I18n.text("Roll"),
                I18n.text("The dice to roll on the table"),
                mLocations.getRoll(),
                new Dice(100, 100, 100),
                new DefaultFormatterFactory(new DiceFormatter(null)),
                (f) -> {
                    mLocations.setRoll((Dice) f.getValue());
                    mAdjustCallback.run();
                });
        if (mFirstField == null) {
            mFirstField = field;
        }
        if (isSubTable()) {
            FontAwesomeButton remove = new FontAwesomeButton("\uf1f8", I18n.text("Remove"), () -> {
                getParent().remove(this);
                mLocations.getOwningLocation().setSubTable(null);
                mAdjustCallback.run();
            });
            wrapper.add(remove);
        }
        add(wrapper, new PrecisionLayoutData().setFillHorizontalAlignment().setGrabHorizontalSpace(true));

        for (HitLocation location : mLocations.getLocations()) {
            add(new HitLocationPanel(location, mAdjustCallback), new PrecisionLayoutData().setGrabHorizontalSpace(true).setFillHorizontalAlignment());
        }
        adjustButtons();
        revalidate();
        EventQueue.invokeLater(this::repaint);
    }

    private static EditorField addField(Container container, String title, String tooltip, Object value, Object protoValue, JFormattedTextField.AbstractFormatterFactory formatter, EditorField.ChangeListener listener) {
        EditorField         field      = new EditorField(formatter, listener, SwingConstants.LEFT, value, protoValue, tooltip);
        PrecisionLayoutData layoutData = new PrecisionLayoutData().setFillHorizontalAlignment();
        if (protoValue == null) {
            layoutData.setGrabHorizontalSpace(true);
        }
        container.add(new Label(title, field), new PrecisionLayoutData().setFillHorizontalAlignment());
        container.add(field, layoutData);
        return field;
    }

    private boolean isSubTable() {
        return mLocations.getOwningLocation() != null;
    }

    public void adjustButtons() {
        Component[] children = getComponents();
        int         count    = children.length;
        for (int i = 1; i < count; i++) {
            ((HitLocationPanel) children[i]).adjustButtons(i == 1, i == count - 1);
        }
    }

    public void adjustForReordering() {
        Component[] children = getComponents();
        int         count    = children.length;
        for (int i = 1; i < count; i++) {
            ((HitLocationPanel) children[i]).adjustButtons(i == 1, i == count - 1);
        }
        mLocations.update();
        repaint();
    }

    public void addHitLocation() {
        HitLocation location = new HitLocation("id", I18n.text("choice name"), I18n.text("table name"), 0, 0, 0, I18n.text("description"));
        mLocations.addLocation(location);
        mLocations.update();
        mAdjustCallback.run();
        add(new HitLocationPanel(location, mAdjustCallback), new PrecisionLayoutData().setGrabHorizontalSpace(true).setFillHorizontalAlignment());
        scrollRectToVisible(new Rectangle(0, getPreferredSize().height - 1, 1, 1));
        ((HitLocationPanel) getComponent(getComponentCount() - 1)).focusIDField();
        adjustButtons();
    }

    public void focusFirstField() {
        mFirstField.requestFocusInWindow();
    }

    @Override
    protected Color getBandingColor(boolean odd) {
        Color color = super.getBandingColor(odd);
        return isSubTable() ? Colors.adjustSaturation(color, 0.05f * countSubTableDepth()) : color;
    }
}
