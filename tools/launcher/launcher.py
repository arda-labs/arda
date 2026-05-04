#!/usr/bin/env python3
"""Arda Launcher — mỗi project 1 tab trong Windows Terminal."""

import os, sys, json, re, subprocess, threading
from pathlib import Path
from dataclasses import dataclass
from typing import Optional

try:
    import tkinter as tk
    from tkinter import ttk, messagebox
except ImportError:
    print("tkinter not available.")
    sys.exit(1)

REPO = Path(__file__).parent.parent.parent.resolve()

# ── Model ───────────────────────────────────────────────────────────

@dataclass
class Project:
    name: str
    path: Path
    cwd: Path
    category: str
    subcategory: str
    command: str
    port: Optional[int] = None

# ── Scanner ─────────────────────────────────────────────────────────

class Scanner:
    def scan(self) -> list[Project]:
        out = []
        self._fe(out); self._go(out); self._java(out)
        return out

    def _fe(self, out):
        f = REPO / "arda/apps/frontend-micro/angular.json"
        if not f.exists(): return
        try: cfg = json.loads(f.read_text("utf-8"))
        except: return
        for name, c in cfg.get("projects", {}).items():
            if c.get("projectType") != "application": continue
            port = None
            try: port = c["architect"]["serve"]["options"]["port"]
            except: pass
            out.append(Project(f"fe/{name}", f.parent / c["root"], f.parent, "fe", "angular-mfe",
                               f"ng serve {name}", port))

    def _go(self, out):
        d = REPO / "arda/apps/backend-go"
        if not d.is_dir(): return
        for child in sorted(d.iterdir()):
            if not child.is_dir(): continue
            gm = child / "go.mod"
            if not gm.exists() or "kratos" not in gm.read_text("utf-8", errors="ignore"): continue
            port = None
            for y in (child / "configs").glob("*.yaml"):
                m = re.search(r'addr:\s*(?:[^:\n]*):(\d{4,5})', y.read_text("utf-8", errors="ignore"))
                if m: port = int(m.group(1)); break
            out.append(Project(f"be/go/{child.name.replace('-service','')}", child, child, "be", "go-kratos",
                               "kratos run", port))

    def _java(self, out):
        d = REPO / "arda/apps/backend-java/settings.gradle.kts"
        if not d.exists(): return
        for line in d.read_text("utf-8", errors="ignore").splitlines():
            if not line.strip().startswith("include("): continue
            svc = line.split("(",1)[1].rsplit(")",1)[0].strip('"').strip("'").lstrip(":")
            sd = d.parent / svc
            if not sd.is_dir(): continue
            port = None
            for f in (sd / "src/main/resources").glob("application*"):
                for l in f.read_text("utf-8", errors="ignore").splitlines():
                    if "server.port" in l:
                        for p in l.replace(":"," ").split():
                            if p.strip().isdigit() and 1024 <= int(p) <= 65535: port = int(p)
            out.append(Project(f"be/java/{svc.replace('-service','')}", sd, d.parent, "be", "java-spring",
                               f'.\\gradlew :{svc}:bootRun -q', port))

# ── Process Manager ─────────────────────────────────────────────────

TAB_COLORS = {
    "angular-mfe": "#89b4fa",
    "go-kratos": "#a6e3a1",
    "java-spring": "#f9e2af",
}

class ProcMgr:
    def __init__(self):
        self._running: dict[str, bool] = {}
        self._scripts: dict[str, Path] = {}

    def start(self, p: Project) -> tuple[bool, str]:
        if self._running.get(p.name):
            return False, f"{p.name} already running."

        title = f"Arda - {p.name}"
        color = TAB_COLORS.get(p.subcategory, "#89b4fa")

        # Strategy 1: wt.exe named-window with PowerShell inside tab
        try:
            wt_cmd = (
                f'wt -w ArdaDev new-tab --title "{title}" --tabColor "{color}" '
                f'-d "{p.cwd}" powershell -NoExit -Command "& {{ {p.command} }}"'
            )
            subprocess.Popen(wt_cmd, shell=True)
            self._running[p.name] = True
            return True, f"Started {p.name}"
        except Exception:
            pass

        # Strategy 2: temp .bat with wt.exe + cmd (fallback)
        try:
            safe_name = p.name.replace("/", "_").replace(":", "")
            script = Path(os.environ["TEMP"]) / f"arda_{safe_name}.bat"
            script.write_text(
                f"@echo off\r\ncd /d {p.cwd}\r\n{p.command}\r\n",
                encoding="utf-8"
            )
            self._scripts[p.name] = script
            subprocess.Popen([
                "wt.exe", "--window", "0",
                "new-tab", "--title", title,
                "cmd.exe", "/k", str(script)
            ])
            self._running[p.name] = True
            return True, f"Started {p.name}"
        except Exception:
            pass

        # Strategy 3: raw cmd.exe window (no WT available)
        try:
            subprocess.Popen(
                ["cmd.exe", "/k", f"cd /d {p.cwd} && {p.command}"],
                creationflags=subprocess.CREATE_NEW_CONSOLE
            )
            self._running[p.name] = True
            return True, f"Started {p.name}"
        except Exception as e:
            return False, str(e)

    def stop(self, name: str) -> tuple[bool, str]:
        if not self._running.get(name):
            return False, f"{name} not running."
        try:
            subprocess.Popen(["taskkill", "/F", "/FI", f'WINDOWTITLE eq Arda - {name}*'],
                             stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
        except: pass
        self._running[name] = False
        if name in self._scripts:
            try: self._scripts[name].unlink()
            except: pass
            del self._scripts[name]
        return True, f"Stopped {name}"

    def is_running(self, name: str) -> bool:
        return bool(self._running.get(name))

    def running_count(self) -> int:
        return sum(1 for v in self._running.values() if v)

# ── UI ──────────────────────────────────────────────────────────────

C = {
    "bg": "#1e1e2e", "card": "#2a2a3e", "hover": "#35354a",
    "fg": "#cdd6f4", "fg2": "#a6adc8",
    "green": "#a6e3a1", "blue": "#89b4fa", "red": "#f38ba8",
    "yellow": "#f9e2af", "mauve": "#cba6f7",
    "border": "#313244",
}
BADGE = {"angular-mfe": (C["blue"], "ANGULAR"), "go-kratos": (C["green"], "GO"), "java-spring": (C["yellow"], "JAVA")}

class App:
    def __init__(self):
        self.root = tk.Tk()
        self.root.title("Arda Launcher")
        self.root.geometry("1000x700")
        self.root.configure(bg=C["bg"])

        self.scanner = Scanner()
        self.pm = ProcMgr()
        self.projects: list[Project] = []
        self.cards: dict[str, dict] = {}

        self._build_ui()
        self.scan()
        self.root.mainloop()

    def _build_ui(self):
        hdr = tk.Frame(self.root, bg=C["bg"])
        hdr.pack(fill="x", padx=32, pady=(20, 0))

        # Logo area
        logo_frame = tk.Frame(hdr, bg=C["bg"])
        logo_frame.pack(side="left")
        logo_dot = tk.Canvas(logo_frame, width=10, height=10, bg=C["bg"],
                             highlightthickness=0)
        logo_dot.create_oval(0, 0, 10, 10, fill=C["blue"], outline="")
        logo_dot.pack(side="left")
        tk.Label(logo_frame, text="Arda Launcher", font=("Segoe UI", 20, "bold"),
                 fg=C["fg"], bg=C["bg"]).pack(side="left", padx=(10, 0))

        # Subtitle
        sub = tk.Label(hdr, text="arda-labs", font=("Segoe UI", 10),
                       fg=C["fg2"], bg=C["bg"])
        sub.pack(side="left", padx=(16, 0), pady=(4, 0))

        # Right side
        right = tk.Frame(hdr, bg=C["bg"])
        right.pack(side="right")

        self.stats = tk.Label(right, font=("Segoe UI", 10, "bold"),
                              fg=C["blue"], bg=C["bg"])
        self.stats.pack(side="right", padx=(12, 0))

        self.scan_btn = self._btn(right, "↻  Rescan", C["blue"], self.scan, padx=14)
        self.scan_btn.pack(side="right")

        # Separator line
        sep = tk.Frame(self.root, height=1, bg=C["border"])
        sep.pack(fill="x", padx=32, pady=(14, 0))

        # Notebook
        style = ttk.Style()
        style.theme_use("clam")
        style.configure("TNotebook", background=C["bg"], borderwidth=0)
        style.configure("TNotebook.Tab", background=C["card"], foreground=C["fg2"],
                        padding=[22, 8], font=("Segoe UI", 10, "bold"))
        style.map("TNotebook.Tab", background=[("selected", C["blue"])],
                   foreground=[("selected", C["bg"])])

        nb = ttk.Notebook(self.root)
        nb.pack(fill="both", expand=True, padx=32, pady=(10, 16))
        self.fe_frame = self._scroll(nb)
        self.be_frame = self._scroll(nb)
        nb.add(self.fe_frame["outer"], text="  Frontend  ")
        nb.add(self.be_frame["outer"], text="  Backend  ")

    def _scroll(self, parent):
        outer = tk.Frame(parent, bg=C["bg"])
        canvas = tk.Canvas(outer, bg=C["bg"], highlightthickness=0, bd=0)
        sb = tk.Scrollbar(outer, orient="vertical", command=canvas.yview,
                          bg=C["card"], troughcolor=C["bg"],
                          activebackground=C["fg2"], bd=0,
                          highlightthickness=0, width=8)
        inner = tk.Frame(canvas, bg=C["bg"])
        inner.bind("<Configure>", lambda e: canvas.configure(scrollregion=canvas.bbox("all")))
        canvas.create_window((0, 0), window=inner, anchor="nw")
        canvas.configure(yscrollcommand=sb.set)
        def _mw(e): canvas.yview_scroll(int(-1*(e.delta/120)), "units")
        canvas.bind("<Enter>", lambda e: canvas.bind_all("<MouseWheel>", _mw, add="+"))
        canvas.bind("<Leave>", lambda e: canvas.unbind_all("<MouseWheel>"))
        canvas.pack(side="left", fill="both", expand=True)
        sb.pack(side="right", fill="y")
        return {"outer": outer, "inner": inner}

    def _btn(self, parent, text, color, cmd, padx=10):
        btn = tk.Button(parent, text=text, font=("Segoe UI", 9, "bold"),
                        fg=C["bg"], bg=color, activebackground=self._lighten(color),
                        relief="flat", padx=padx, pady=5, cursor="hand2", command=cmd,
                        bd=0, highlightthickness=0)
        btn.bind("<Enter>", lambda e: btn.configure(bg=self._lighten(color, 0.32)))
        btn.bind("<Leave>", lambda e: btn.configure(bg=color))
        return btn

    @staticmethod
    def _lighten(c, f=0.18):
        r = tuple(int(c.lstrip("#")[i:i+2],16) for i in (0,2,4))
        r2 = tuple(min(255,int(v+(255-v)*f)) for v in r)
        return f"#{r2[0]:02x}{r2[1]:02x}{r2[2]:02x}"

    def scan(self):
        self.scan_btn.configure(text="⏳", state="disabled")
        self.root.update()
        self.projects = self.scanner.scan()
        self._render()
        self.scan_btn.configure(text="Rescan", state="normal")

    def _render(self):
        for f in (self.fe_frame["inner"], self.be_frame["inner"]):
            for w in f.winfo_children(): w.destroy()
        self.cards.clear()
        for p in self.projects:
            parent = self.fe_frame["inner"] if p.category == "fe" else self.be_frame["inner"]
            self._card(parent, p)
        self._update_stats()

    def _card(self, parent, p):
        card = tk.Frame(parent, bg=C["card"], bd=0, highlightthickness=1,
                        highlightbackground=C["border"], highlightcolor=C["border"])
        card.pack(fill="x", padx=6, pady=4, ipady=0)

        badge_clr, badge_txt = BADGE.get(p.subcategory, (C["mauve"], p.subcategory.upper()))

        # Accent bar — 4px colored strip on the left
        accent = tk.Frame(card, width=4, bg=badge_clr, bd=0)
        accent.pack(side="left", fill="y")
        accent.pack_propagate(False)

        body = tk.Frame(card, bg=C["card"], bd=0)
        body.pack(side="left", fill="x", expand=True, padx=(10, 6), pady=8)

        # Top row: badge + name + port + status dot (right-aligned)
        top = tk.Frame(body, bg=C["card"])
        top.pack(fill="x")

        # Badge label
        badge = tk.Label(top, text=badge_txt, font=("Segoe UI", 7, "bold"),
                         fg=C["bg"], bg=badge_clr, padx=6, pady=1)
        badge.pack(side="left", padx=(0, 6))

        # Project name
        name_lbl = tk.Label(top, text=p.name, font=("Segoe UI", 12, "bold"),
                            fg=C["fg"], bg=C["card"])
        name_lbl.pack(side="left")

        # Port chip
        if p.port:
            port_lbl = tk.Label(top, text=f":{p.port}", font=("Segoe UI", 8, "bold"),
                                fg=C["bg"], bg=C["blue"], padx=5, pady=0)
            port_lbl.pack(side="left", padx=(6, 0))

        # Status dot
        st = tk.Label(top, text="●", font=("Segoe UI", 14), fg=C["fg2"], bg=C["card"])
        st.pack(side="right", padx=(0, 2))

        # Bottom row: path + command
        try: rp = p.path.relative_to(REPO)
        except: rp = p.path
        path_lbl = tk.Label(body, text=f"{rp}", font=("Segoe UI", 8),
                            fg=C["fg2"], bg=C["card"], anchor="w")
        path_lbl.pack(fill="x")
        cmd_lbl = tk.Label(body, text=f"$ {p.command}", font=("Segoe UI", 9, "italic"),
                           fg=C["mauve"], bg=C["card"], anchor="w")
        cmd_lbl.pack(fill="x")

        # Buttons column
        act = tk.Frame(card, bg=C["card"])
        act.pack(side="right", fill="y", padx=(0, 8), pady=8)

        start_btn = self._card_btn(act, "▶", "Start", C["green"],
                                   lambda n=p.name: self._start(n))
        start_btn.pack(side="top", pady=(0, 2))

        stop_btn = self._card_btn(act, "■", "Stop", C["red"],
                                  lambda n=p.name: self._stop(n))
        stop_btn.pack(side="top")

        # Hover
        def _on_enter(e):
            card.configure(bg=C["hover"], highlightbackground=C["hover"])
            for w in (body, top, name_lbl, st, path_lbl, cmd_lbl, act, badge):
                if w and w.winfo_exists():
                    try: w.configure(bg=C["hover"])
                    except: pass
        def _on_leave(e):
            card.configure(bg=C["card"], highlightbackground=C["border"])
            for w in (body, top, name_lbl, st, path_lbl, cmd_lbl, act, badge):
                if w and w.winfo_exists():
                    try: w.configure(bg=C["card"])
                    except: pass
        card.bind("<Enter>", _on_enter)
        card.bind("<Leave>", _on_leave)
        # Forward enter/leave to card via body events too
        body.bind("<Enter>", _on_enter)
        body.bind("<Leave>", _on_leave)

        self.cards[p.name] = {"status": st, "card": card}

    def _card_btn(self, parent, icon, tip, color, cmd):
        btn = tk.Button(parent, text=icon, font=("Segoe UI", 10),
                        fg=color, bg=C["card"], activeforeground=C["bg"],
                        activebackground=color,
                        relief="flat", padx=6, pady=2, cursor="hand2",
                        bd=0, highlightthickness=0,
                        command=cmd)
        self._attach_hint(btn, tip)
        btn.bind("<Enter>", lambda e: btn.configure(fg=C["bg"], bg=color))
        btn.bind("<Leave>", lambda e: btn.configure(fg=color, bg=C["card"]))
        return btn

    @staticmethod
    def _attach_hint(widget, text):
        tip = None
        def _show(e):
            nonlocal tip
            if tip: return
            import tkinter as tk
            tip = tk.Toplevel(widget)
            tip.wm_overrideredirect(True)
            tip.wm_geometry(f"+{e.x_root+12}+{e.y_root-6}")
            lbl = tk.Label(tip, text=text, font=("Segoe UI", 8),
                           fg=C["fg"], bg=C["card"],
                           relief="solid", bd=1, padx=6, pady=2)
            lbl.pack()
        def _move(e):
            nonlocal tip
            if tip: tip.wm_geometry(f"+{e.x_root+12}+{e.y_root-6}")
        def _hide(e):
            nonlocal tip
            if tip: tip.destroy(); tip = None
        widget.bind("<Enter>", _show)
        widget.bind("<Motion>", _move)
        widget.bind("<Leave>", _hide)

    def _start(self, name):
        p = next((x for x in self.projects if x.name == name), None)
        if not p: return
        # Show pending state immediately
        d = self.cards.get(name)
        if d: d["status"].configure(fg=C["yellow"])
        threading.Thread(target=self._start_task, args=(p,), daemon=True).start()

    def _start_task(self, p):
        ok, msg = self.pm.start(p)
        if not ok: self.root.after(0, lambda: messagebox.showwarning("", msg))
        self.root.after(0, lambda: self._refresh(p.name))

    def _stop(self, name):
        d = self.cards.get(name)
        if d: d["status"].configure(fg=C["yellow"])
        threading.Thread(target=self._stop_task, args=(name,), daemon=True).start()

    def _stop_task(self, name):
        ok, msg = self.pm.stop(name)
        if not ok: self.root.after(0, lambda: messagebox.showwarning("", msg))
        self.root.after(0, lambda: self._refresh(name))

    def _refresh(self, name):
        d = self.cards.get(name)
        if d:
            d["status"].configure(fg=C["green"] if self.pm.is_running(name) else C["fg2"])
        self._update_stats()

    def _update_stats(self):
        self.stats.configure(text=f"{self.pm.running_count()}/{len(self.projects)} running")

if __name__ == "__main__":
    App()
