package ac

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const musicRoot = "/Net/music/alac"

type letterBucket struct {
	label string
	start string // "#" or first letter in bucket
}

var letterBuckets = []letterBucket{
	{label: "#–C", start: "#"},
	{label: "D–G", start: "D"},
	{label: "H–K", start: "H"},
	{label: "L–O", start: "L"},
	{label: "P–S", start: "P"},
	{label: "T–W", start: "T"},
	{label: "X–Z", start: "X"},
}

type musicItem struct {
	label  string
	path   string
	isDir  bool
	action string // "", "back", "playall", "shuffle"
	bounds image.Rectangle
}

type musicBrowser struct {
	mu          sync.Mutex
	artists     []string // cached top-level directories
	currentPath string   // "" = artist root, else full dir path
	items       []musicItem
	page        int
	pageUp      *sprite
	pageDn      *sprite
	bucketRects map[string]image.Rectangle

	queue      []string
	queueIdx   int
	nowPlaying string
}

func musicDialog(g *Game) {
	chirp(g)
	if g.music == nil {
		g.music = &musicBrowser{
			bucketRects: make(map[string]image.Rectangle),
		}
	}
	g.music.ensureLoaded(g)
	g.state = inMusic
}

func (m *musicBrowser) ensureLoaded(g *Game) {
	if m.pageUp == nil {
		m.pageUp = getSprite("Up", "", nil)
		m.pageUp.scale = 1.4
		m.pageDn = getSprite("Dn", "", nil)
		m.pageDn.scale = 1.4
	}
	if len(m.artists) == 0 {
		entries, err := os.ReadDir(musicRoot)
		if err != nil {
			g.debug(err.Error())
			return
		}
		for _, e := range entries {
			if e.IsDir() && !strings.HasPrefix(e.Name(), ".") {
				m.artists = append(m.artists, e.Name())
			}
		}
		slices.SortFunc(m.artists, func(a, b string) int {
			ca, cb := leadingBucketChar(a), leadingBucketChar(b)
			if ca != cb {
				if ca == '#' {
					return -1
				}
				if cb == '#' {
					return 1
				}
				if ca < cb {
					return -1
				}
				return 1
			}
			return strings.Compare(strings.ToLower(a), strings.ToLower(b))
		})
	}
	if len(m.items) == 0 {
		m.openRoot()
	}
}

func (m *musicBrowser) openRoot() {
	m.currentPath = ""
	m.page = 0
	m.items = make([]musicItem, 0, len(m.artists)+1)
	m.items = append(m.items, musicItem{
		label:  "Shuffle Library",
		action: "shuffle",
	})
	for _, a := range m.artists {
		m.items = append(m.items, musicItem{
			label: a,
			path:  filepath.Join(musicRoot, a),
			isDir: true,
		})
	}
}

func (m *musicBrowser) openDir(g *Game, dirPath string) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		g.debug(err.Error())
		return
	}

	m.currentPath = dirPath
	m.page = 0
	m.items = []musicItem{
		{label: ".. Back", action: "back"},
		{label: "Play All", path: dirPath, action: "playall"},
	}

	var dirs, files []musicItem
	for _, e := range entries {
		name := e.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}
		full := filepath.Join(dirPath, name)
		if e.IsDir() {
			dirs = append(dirs, musicItem{label: name, path: full, isDir: true})
		} else if isAudioFile(name) {
			files = append(files, musicItem{label: trimAudioExt(name), path: full, isDir: false})
		}
	}

	slices.SortFunc(dirs, func(a, b musicItem) int {
		return strings.Compare(strings.ToLower(a.label), strings.ToLower(b.label))
	})
	slices.SortFunc(files, func(a, b musicItem) int {
		return strings.Compare(strings.ToLower(a.label), strings.ToLower(b.label))
	})

	m.items = append(m.items, dirs...)
	m.items = append(m.items, files...)
}

func isAudioFile(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	switch ext {
	case ".m4a", ".alac", ".flac", ".mp3", ".ogg", ".opus", ".wav":
		return true
	}
	return false
}

func trimAudioExt(name string) string {
	ext := filepath.Ext(name)
	return strings.TrimSuffix(name, ext)
}

func collectTracksRecursive(dir string) []string {
	var tracks []string
	_ = filepath.WalkDir(dir, func(p string, d os.DirEntry, err error) error {
		if err != nil || d == nil {
			return nil
		}
		if d.IsDir() && strings.HasPrefix(d.Name(), ".") {
			return filepath.SkipDir
		}
		if !d.IsDir() && isAudioFile(d.Name()) {
			tracks = append(tracks, p)
		}
		return nil
	})
	slices.Sort(tracks)
	return tracks
}

func (g *Game) startMusicQueue(tracks []string, startIdx int) {
	if len(tracks) == 0 {
		return
	}
	if startIdx < 0 || startIdx >= len(tracks) {
		startIdx = 0
	}

	g.stopPlayer()

	g.music.mu.Lock()
	g.music.queue = tracks
	g.music.queueIdx = startIdx
	g.music.mu.Unlock()

	g.playQueueIndex(startIdx)
}

func (g *Game) playQueueIndex(idx int) {
	g.music.mu.Lock()
	if idx < 0 || idx >= len(g.music.queue) {
		g.music.queue = nil
		g.music.nowPlaying = ""
		g.music.mu.Unlock()
		return
	}
	g.music.queueIdx = idx
	trackPath := g.music.queue[idx]
	artist := filepath.Base(filepath.Dir(filepath.Dir(trackPath)))
	title := trimAudioExt(filepath.Base(trackPath))
	if artist == "alac" || artist == "." || artist == "/" {
		g.music.nowPlaying = title
	} else {
		g.music.nowPlaying = fmt.Sprintf("%s — %s", artist, title)
	}
	g.music.mu.Unlock()

	g.playExternalWithCallback(trackPath, func() {
		g.nextMusicTrack()
	})
}

func (g *Game) nextMusicTrack() {
	if g.music == nil {
		return
	}
	g.music.mu.Lock()
	hasQueue := len(g.music.queue) > 0
	nextIdx := g.music.queueIdx + 1
	g.music.mu.Unlock()

	if !hasQueue {
		return
	}
	if nextIdx >= len(g.music.queue) {
		g.stopPlayer()
		return
	}
	g.playQueueIndex(nextIdx)
}

func (g *Game) hasActiveMusicQueue() bool {
	if g.music == nil {
		return false
	}
	g.music.mu.Lock()
	defer g.music.mu.Unlock()
	return len(g.music.queue) > 0
}

func (g *Game) currentNowPlaying() string {
	if g.music == nil {
		return ""
	}
	g.music.mu.Lock()
	defer g.music.mu.Unlock()
	return g.music.nowPlaying
}

const musicRowsPerPage = 5

func (g *Game) drawMusicDialog(screen *ebiten.Image) {
	g.drawModal(screen)
	if g.music == nil {
		return
	}

	mx, my, mw, mh := modalBounds()
	closeBtn := modalbuttons[0]
	closeW := int(float64(closeBtn.image.Bounds().Max.X) * closeBtn.scale)

	innerX := mx + modalPadding
	innerY := my + modalPadding
	rightLimit := mx + mw - closeW - modalPadding*2

	// Header breadcrumb / status
	header := "Music Library"
	if g.music.currentPath != "" {
		rel, _ := filepath.Rel(musicRoot, g.music.currentPath)
		header = rel
	}
	drawShadowedText(screen, truncateLabel(header, 38), controlfont, float64(innerX), float64(innerY)+4, textPrimary)

	// Finger-friendly letter-range jump bar at root level (~85x44px buttons)
	listTopY := innerY + 32
	if g.music.currentPath == "" {
		stripW := rightLimit - innerX
		gap := 6
		cellW := (stripW - gap*(len(letterBuckets)-1)) / len(letterBuckets)
		cellH := 44
		for i, b := range letterBuckets {
			lx := innerX + i*(cellW+gap)
			rect := image.Rect(lx, listTopY, lx+cellW, listTopY+cellH)
			g.music.bucketRects[b.start] = rect
			drawRoundedBox(screen, float32(rect.Min.X), float32(rect.Min.Y), float32(rect.Dx()), float32(rect.Dy()), 10, cardSlate, bordergrey, 1.0)
			tw, th := text.Measure(b.label, controlfont, 0)
			tx := float64(rect.Min.X) + (float64(rect.Dx())-tw)/2
			ty := float64(rect.Min.Y) + (float64(rect.Dy())-th)/2
			drawShadowedText(screen, b.label, controlfont, tx, ty, textPrimary)
		}
		listTopY += cellH + 10
	}

	// Paging controls column on the right of the list (below Close button)
	navColW := 64
	listW := rightLimit - innerX - navColW - 10
	listH := (my + mh - modalPadding) - listTopY
	rowGap := 8
	rowH := (listH - rowGap*(musicRowsPerPage-1)) / musicRowsPerPage

	totalPages := (len(g.music.items) + musicRowsPerPage - 1) / musicRowsPerPage
	if totalPages < 1 {
		totalPages = 1
	}
	if g.music.page >= totalPages {
		g.music.page = totalPages - 1
	}
	if g.music.page < 0 {
		g.music.page = 0
	}

	// Clear old item hitboxes before populating current page
	for i := range g.music.items {
		g.music.items[i].bounds = image.Rectangle{}
	}

	start := g.music.page * musicRowsPerPage
	end := start + musicRowsPerPage
	if end > len(g.music.items) {
		end = len(g.music.items)
	}

	for rowIdx, idx := range makeRange(start, end) {
		item := &g.music.items[idx]
		ry := listTopY + rowIdx*(rowH+rowGap)
		item.bounds = image.Rect(innerX, ry, innerX+listW, ry+rowH)

		fill := cardSlate
		border := bordergrey
		bw := float32(1.0)
		labelCol := textPrimary

		if item.action == "playall" || item.action == "shuffle" {
			fill = pinkTint
			border = pink
			bw = 1.5
		} else if item.action == "back" {
			labelCol = textMuted
		}

		drawRoundedBox(screen, float32(innerX), float32(ry), float32(listW), float32(rowH), 12, fill, border, bw)

		prefix := ""
		if item.isDir {
			prefix = "▸ "
		}
		_, th := text.Measure(item.label, controlfont, 0)
		ty := float64(ry) + (float64(rowH)-th)/2
		drawShadowedText(screen, prefix+truncateLabel(item.label, 42), controlfont, float64(innerX+14), ty, labelCol)
	}

	// Page Up / Page indicator / Page Down
	navX := innerX + listW + 10
	upB := g.music.pageUp.image.Bounds()
	upW := int(float64(upB.Max.X) * g.music.pageUp.scale)
	upH := int(float64(upB.Max.Y) * g.music.pageUp.scale)

	g.music.pageUp.setLocation(navX+(navColW-upW)/2, listTopY+4)
	g.music.pageUp.draw(screen)

	pageLabel := fmt.Sprintf("%d/%d", g.music.page+1, totalPages)
	pw, ph := text.Measure(pageLabel, controlfont, 0)
	drawShadowedText(screen, pageLabel, controlfont, float64(navX)+(float64(navColW)-pw)/2, float64(listTopY)+float64(listH)/2-ph/2, textMuted)

	g.music.pageDn.setLocation(navX+(navColW-upW)/2, listTopY+listH-upH-4)
	g.music.pageDn.draw(screen)
}

func (g *Game) handleMusicClick(x, y int) {
	if g.music == nil {
		return
	}

	// Page Up / Dn
	if g.music.pageUp != nil && g.music.pageUp.in(x, y) {
		if g.music.page > 0 {
			g.music.page--
			tick(g)
		}
		return
	}
	if g.music.pageDn != nil && g.music.pageDn.in(x, y) {
		totalPages := (len(g.music.items) + musicRowsPerPage - 1) / musicRowsPerPage
		if g.music.page+1 < totalPages {
			g.music.page++
			tick(g)
		}
		return
	}

	// Letter-bucket quick jump when at artist root
	if g.music.currentPath == "" {
		for startLetter, rect := range g.music.bucketRects {
			if x >= rect.Min.X && x <= rect.Max.X && y >= rect.Min.Y && y <= rect.Max.Y {
				g.music.jumpToLetter(startLetter)
				tick(g)
				return
			}
		}
	}

	// List item taps
	for i := range g.music.items {
		item := g.music.items[i]
		if item.bounds.Empty() {
			continue
		}
		if x >= item.bounds.Min.X && x <= item.bounds.Max.X && y >= item.bounds.Min.Y && y <= item.bounds.Max.Y {
			switch item.action {
			case "back":
				tick(g)
				parent := filepath.Dir(g.music.currentPath)
				if parent == musicRoot || len(parent) <= len(musicRoot) {
					g.music.openRoot()
				} else {
					g.music.openDir(g, parent)
				}
			case "playall":
				chirp(g)
				tracks := collectTracksRecursive(item.path)
				g.startMusicQueue(tracks, 0)
				g.state = inNormal
			case "shuffle":
				chirp(g)
				g.shuffleRandomArtist(g)
				g.state = inNormal
			default:
				if item.isDir {
					tick(g)
					g.music.openDir(g, item.path)
				} else {
					chirp(g)
					var tracks []string
					startIdx := 0
					for _, sib := range g.music.items {
						if !sib.isDir && sib.action == "" && sib.path != "" {
							if sib.path == item.path {
								startIdx = len(tracks)
							}
							tracks = append(tracks, sib.path)
						}
					}
					g.startMusicQueue(tracks, startIdx)
					g.state = inNormal
				}
			}
			return
		}
	}
}

func (g *Game) shuffleRandomArtist(game *Game) {
	if len(g.music.artists) == 0 {
		return
	}
	var pool []string
	for i := 0; i < 15; i++ {
		artist := g.music.artists[rand.Intn(len(g.music.artists))]
		pool = append(pool, collectTracksRecursive(filepath.Join(musicRoot, artist))...)
	}
	rand.Shuffle(len(pool), func(i, j int) { pool[i], pool[j] = pool[j], pool[i] })
	game.startMusicQueue(pool, 0)
}

func (m *musicBrowser) jumpToLetter(letter string) {
	if len(letter) == 0 {
		return
	}
	target := letter[0]
	for idx, item := range m.items {
		if item.action != "" || len(item.label) == 0 {
			continue
		}
		c := leadingBucketChar(item.label)
		if target == '#' {
			if c == '#' {
				m.page = idx / musicRowsPerPage
				return
			}
		} else if c != '#' && c >= target {
			m.page = idx / musicRowsPerPage
			return
		}
	}
}

func drawShadowedText(dst *ebiten.Image, str string, face *text.GoTextFace, x, y float64, col color.RGBA) {
	shadowOp := &text.DrawOptions{}
	shadowOp.GeoM.Translate(x, y+1)
	shadowOp.ColorScale.ScaleWithColor(textShadow)
	text.Draw(dst, str, face, shadowOp)

	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(col)
	text.Draw(dst, str, face, op)
}

func truncateLabel(s string, maxRunes int) string {
	r := []rune(s)
	if len(r) <= maxRunes {
		return s
	}
	return string(r[:maxRunes-1]) + "…"
}

func makeRange(start, end int) []int {
	out := make([]int, 0, end-start)
	for i := start; i < end; i++ {
		out = append(out, i)
	}
	return out
}

func leadingBucketChar(s string) byte {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return '#'
	}
	c := s[0]
	if c >= 'a' && c <= 'z' {
		c = c - 'a' + 'A'
	}
	if c >= 'A' && c <= 'Z' {
		return c
	}
	return '#'
}
