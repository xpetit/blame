### Blame!

```
27.0% [████████▎   ] librewolf     2.2 GB [█████▋      ] librewolf
 9.2% [██▉         ] codium        2.1 GB [█████▍      ] codium
 1.9% [▋           ] Xorg           85 MB [▎           ] sublime_merge
 0.2% [▏           ] xfce4-panel    58 MB [▏           ] Xorg
 0.1% [            ] xfwm4          22 MB [            ] xfwm4
 0.1% [            ] pulseaudio     20 MB [            ] xfdesktop
                                    16 MB [            ] statusbar
                                    16 MB [            ] xfce4-session
                                    13 MB [            ] xfce4-appfinder
                                   9.9 MB [            ] thunar
                                   9.4 MB [            ] pulseaudio
                                   8.5 MB [            ] wrapper-2.0 (panel-8-pulseau)
```

Shows the Linux processes that use more than 0.1% of your resources (CPU & RAM).
Groups processes belonging to a window together (X11 only).

The graphs show the resource usage compared to other processes.
A half-full bar means that this process uses as much as all the other processes combined.

### Usage

`wmctrl` is required.

```
go run github.com/xpetit/blame@latest
```

### TODO

- Support
  - Multiple monitors
  - Multiple workspaces
- User-provided process parents list (name or PID)
- Improve code (simplify & add comments)
