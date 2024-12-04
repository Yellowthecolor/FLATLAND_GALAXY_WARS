package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Audio struct {
	currentMusic Music
	menuMusic    Music
	gameMusic    Music
	winMusic     Music
	loseMusic    Music

	collectSound SoundSlice
	fiveSound    SoundSlice
	impactSound  SoundSlice
	spawnSound   SoundSlice
	warningSound SoundSlice
	pewSound     SoundSlice

	musicChanged bool
}

type Music struct {
	music  rl.Music
	volume float32
}

type SoundSlice struct {
	soundSlice   []rl.Sound
	maxSounds    int
	currentSound int
	volume       float32
}

func (sc *Scene) LoadMusic() {
	sc.menuMusic.music = loadMusicFromEmbed("Audio/Music/menu.wav")
	sc.menuMusic.volume = 0.5

	sc.gameMusic.music = loadMusicFromEmbed("Audio/Music/game.wav")
	sc.gameMusic.volume = 0.1

	sc.winMusic.music = loadMusicFromEmbed("Audio/Music/win.wav")
	sc.winMusic.volume = 0.1

	sc.loseMusic.music = loadMusicFromEmbed("Audio/Music/lose.wav")
	sc.loseMusic.volume = 0.1
}

func loadMusicFromEmbed(path string) rl.Music {
	data, err := audioFiles.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("failed to read embedded music file %s: %w", path, err))
	}
	music := rl.LoadMusicStreamFromMemory(".wav", data, int32(len(data)))
	return music
}

func (sc *Scene) LoadSounds() {
	sc.collectSound = NewSoundSlice("collect", 0.1, false)
	sc.fiveSound = NewSoundSlice("five", 0.5, false)
	sc.impactSound = NewSoundSlice("impact", 0.2, true)
	sc.spawnSound = NewSoundSlice("spawn", 0.2, false)
	sc.warningSound = NewSoundSlice("warning", 0.7, false)
	sc.pewSound = NewSoundSlice("pew", 0.2, true)
}

func NewSoundSlice(soundName string, volume float32, overlap bool) SoundSlice {
	newSoundSlice := SoundSlice{
		soundSlice:   make([]rl.Sound, 0),
		maxSounds:    5,
		currentSound: 0,
		volume:       volume,
	}
	newSoundSlice.InitSound(soundName, overlap)
	return newSoundSlice
}

func (ss *SoundSlice) InitSound(soundName string, overlap bool) {
	sound := loadSoundFromEmbed("Audio/Sounds/" + soundName + ".wav")
	if overlap {
		ss.AddRepeatSoundToSlice(sound)
	} else {
		ss.AddSoundToSlice(sound)
	}
}

func loadSoundFromEmbed(path string) rl.Sound {
	data, err := audioFiles.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("failed to read embedded sound file %s: %w", path, err))
	}
	wave := rl.LoadWaveFromMemory(".wav", data, int32(len(data)))
	sound := rl.LoadSoundFromWave(wave)
	return sound
}

func (ss *SoundSlice) AddRepeatSoundToSlice(newSound rl.Sound) {
	ss.soundSlice = append(ss.soundSlice, newSound)
	for i := 1; i < ss.maxSounds; i++ {
		ss.soundSlice = append(ss.soundSlice, rl.LoadSoundAlias(ss.soundSlice[0]))
	}
}

func (ss *SoundSlice) AddSoundToSlice(newSound rl.Sound) {
	ss.soundSlice = append(ss.soundSlice, newSound)
}

func PlaySoundOverlap(ss *SoundSlice) {
	rl.SetSoundVolume(ss.soundSlice[ss.currentSound], ss.volume)
	rl.PlaySound(ss.soundSlice[ss.currentSound])
	ss.currentSound++
	if ss.currentSound >= len(ss.soundSlice) {
		ss.currentSound = 0
	}
}

func PlaySoundStandAlone(ss *SoundSlice) {
	rl.SetSoundVolume(ss.soundSlice[ss.currentSound], ss.volume)
	rl.PlaySound(ss.soundSlice[ss.currentSound])
}

func (sc *Scene) PlayMusic(nextTrack Music) {
	if sc.currentMusic != nextTrack {
		rl.StopMusicStream(sc.currentMusic.music)
	}

	sc.currentMusic = nextTrack
	rl.SetMusicVolume(nextTrack.music, nextTrack.volume)
	rl.PlayMusicStream(nextTrack.music)
}

func (sc *Scene) LoadAudio() {
	sc.LoadMusic()
	sc.LoadSounds()
}
