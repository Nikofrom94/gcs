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

package com.trollworks.gcs.settings;

import com.trollworks.gcs.menu.file.CloseHandler;
import com.trollworks.gcs.ui.Theme;
import com.trollworks.gcs.ui.ThemeColor;
import com.trollworks.gcs.ui.ThemeFont;
import com.trollworks.gcs.ui.UIUtilities;
import com.trollworks.gcs.ui.layout.PrecisionLayout;
import com.trollworks.gcs.ui.layout.PrecisionLayoutAlignment;
import com.trollworks.gcs.ui.layout.PrecisionLayoutData;
import com.trollworks.gcs.ui.widget.BaseWindow;
import com.trollworks.gcs.ui.widget.ColorWell;
import com.trollworks.gcs.ui.widget.FontAwesomeButton;
import com.trollworks.gcs.ui.widget.FontPanel;
import com.trollworks.gcs.ui.widget.Separator;
import com.trollworks.gcs.ui.widget.Label;
import com.trollworks.gcs.ui.widget.Panel;
import com.trollworks.gcs.ui.widget.ScrollPanel;
import com.trollworks.gcs.ui.widget.WindowUtils;
import com.trollworks.gcs.utility.I18n;

import java.awt.BorderLayout;
import java.awt.Color;
import java.awt.Container;
import java.awt.Dimension;
import java.awt.Font;
import java.awt.event.WindowEvent;
import java.util.ArrayList;
import java.util.List;

/** A window for editing theme settings. */
public final class ThemeSettingsWindow extends BaseWindow implements CloseHandler {
    private static ThemeSettingsWindow INSTANCE;
    private        List<FontTracker>   mFontPanels;
    private        List<ColorTracker>  mColorWells;
    private        FontAwesomeButton   mResetFontsButton;
    private        FontAwesomeButton   mResetColorsButton;
    private        boolean             mIgnore;

    /** Displays the theme settings window. */
    public static void display() {
        if (!UIUtilities.inModalState()) {
            ThemeSettingsWindow wnd;
            synchronized (ThemeSettingsWindow.class) {
                if (INSTANCE == null) {
                    INSTANCE = new ThemeSettingsWindow();
                }
                wnd = INSTANCE;
            }
            wnd.setVisible(true);
        }
    }

    private ThemeSettingsWindow() {
        super(I18n.text("Theme Settings"));
        Panel panel = new Panel(new PrecisionLayout().setMargins(10));

        mResetFontsButton = addHeader(panel, I18n.text("Fonts"), 0, this::resetFonts);
        Panel wrapper = new Panel(new PrecisionLayout().setColumns(2), false);
        mFontPanels = new ArrayList<>();
        for (ThemeFont font : ThemeFont.ALL) {
            if (font.isEditable()) {
                FontTracker tracker = new FontTracker(font);
                wrapper.add(new Label(font.toString(), tracker), new PrecisionLayoutData().setFillHorizontalAlignment());
                wrapper.add(tracker);
                mFontPanels.add(tracker);
            }
        }
        panel.add(wrapper);

        mResetColorsButton = addHeader(panel, I18n.text("Colors"), 16, this::resetColors);
        int cols = 8;
        wrapper = new Panel(new PrecisionLayout().setColumns(cols), false);
        mColorWells = new ArrayList<>();
        int max = ThemeColor.ALL.size();
        cols /= 2;
        int maxPerCol  = max / cols;
        int excess     = max % (maxPerCol * cols);
        int iterations = maxPerCol;
        if (excess != 0) {
            iterations++;
        }
        for (int i = 0; i < iterations; i++) {
            addColorTracker(wrapper, ThemeColor.ALL.get(i), 0);
            int index = i;
            for (int j = 1; j < ((i == maxPerCol) ? excess : cols); j++) {
                index += maxPerCol;
                if (j - 1 < excess) {
                    index++;
                }
                if (index < max) {
                    addColorTracker(wrapper, ThemeColor.ALL.get(index), 8);
                }
            }
        }
        panel.add(wrapper);

        getContentPane().add(new ScrollPanel(panel), BorderLayout.CENTER);
        adjustResetButtons();
        establishSizing();
        WindowUtils.packAndCenterWindowOn(this, null);
    }

    @Override
    public void establishSizing() {
        pack();
        int width = getSize().width;
        setMinimumSize(new Dimension(width, 200));
        setMaximumSize(new Dimension(width, getPreferredSize().height));
    }

    private static FontAwesomeButton addHeader(Container parent, String text, int topMargin, Runnable reset) {
        Panel header = new Panel(new PrecisionLayout().setColumns(2).setMargins(0));
        Label label  = new Label(text);
        label.setThemeFont(ThemeFont.HEADER);
        header.add(label);
        FontAwesomeButton resetButton = new FontAwesomeButton("\uf011", I18n.text("Reset to Factory Defaults"), reset);
        header.add(resetButton, new PrecisionLayoutData().setGrabHorizontalSpace(true).setHorizontalAlignment(PrecisionLayoutAlignment.END));
        header.add(new Separator(), new PrecisionLayoutData().setFillHorizontalAlignment().setGrabHorizontalSpace(true).setHorizontalSpan(2));
        parent.add(header, new PrecisionLayoutData().setFillHorizontalAlignment().setGrabHorizontalSpace(true).setTopMargin(topMargin));
        return resetButton;
    }

    private void addColorTracker(Container parent, ThemeColor color, int leftMargin) {
        ColorTracker tracker = new ColorTracker(color);
        mColorWells.add(tracker);
        parent.add(new Label(color.toString(), tracker), new PrecisionLayoutData().setFillHorizontalAlignment().setLeftMargin(leftMargin));
        parent.add(tracker, new PrecisionLayoutData().setLeftMargin(4));
    }

    private void resetFonts() {
        mIgnore = true;
        for (FontTracker tracker : mFontPanels) {
            tracker.reset();
        }
        mIgnore = false;
        BaseWindow.forceRevalidateAndRepaint();
        adjustResetButtons();
    }

    private void resetColors() {
        mIgnore = true;
        for (ColorTracker tracker : mColorWells) {
            tracker.reset();
        }
        mIgnore = false;
        Theme.repaint();
        adjustResetButtons();
    }

    private void adjustResetButtons() {
        boolean enabled = false;
        for (ThemeFont font : ThemeFont.ALL) {
            if (font.isEditable() && !font.getFont().equals(Theme.DEFAULT.getFont(font.getIndex()))) {
                enabled = true;
                break;
            }
        }
        mResetFontsButton.setEnabled(enabled);

        enabled = false;
        for (ThemeColor color : ThemeColor.ALL) {
            if (color.getRGB() != Theme.DEFAULT.getColor(color.getIndex()).getRGB()) {
                enabled = true;
                break;
            }
        }
        mResetColorsButton.setEnabled(enabled);
    }

    @Override
    public boolean mayAttemptClose() {
        return true;
    }

    @Override
    public boolean attemptClose() {
        windowClosing(new WindowEvent(this, WindowEvent.WINDOW_CLOSING));
        return true;
    }


    @Override
    public void dispose() {
        synchronized (ThemeSettingsWindow.class) {
            INSTANCE = null;
        }
        super.dispose();
    }

    private class FontTracker extends FontPanel {
        private int mIndex;

        FontTracker(ThemeFont font) {
            super(font.getFont());
            mIndex = font.getIndex();
            addActionListener((evt) -> {
                if (!mIgnore) {
                    Theme.current().setFont(mIndex, getCurrentFont());
                    adjustResetButtons();
                    BaseWindow.forceRevalidateAndRepaint();
                }
            });
        }

        void reset() {
            Font font = Theme.DEFAULT.getFont(mIndex);
            setCurrentFont(font);
            Theme.current().setFont(mIndex, font);
        }
    }

    private class ColorTracker extends ColorWell implements ColorWell.ColorChangedListener {
        private int mIndex;

        ColorTracker(ThemeColor color) {
            super(new Color(color.getRGB(), true), null);
            mIndex = color.getIndex();
            setColorChangedListener(this);
        }

        void reset() {
            Color color = Theme.DEFAULT.getColor(mIndex);
            setWellColor(color);
            Theme.current().setColor(mIndex, color);
        }

        @Override
        public void colorChanged(Color color) {
            if (!mIgnore) {
                Theme.current().setColor(mIndex, color);
                adjustResetButtons();
                Theme.repaint();
            }
        }
    }
}
