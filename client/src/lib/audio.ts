type AudioTrack = "city-building" | "road-building" | "settlement-building";

class SettlersAudio {
  private currentPlayingTrack: AudioTrack | null = null;
  private audios: Record<AudioTrack, HTMLAudioElement> = {
    "city-building": new Audio("/static/city-building.mp3"),
    "road-building": new Audio("/static/road-building.mp3"),
    "settlement-building": new Audio("/static/settlement-building.mp3"),
  };

  private onEndedCallback() {
    this.currentPlayingTrack = null;
  }

  stopAudio() {
    if (this.currentPlayingTrack) {
      const currentPlayingAudio = this.audios[this.currentPlayingTrack];
      currentPlayingAudio.pause();
      currentPlayingAudio.currentTime = 0;
      currentPlayingAudio.removeEventListener("ended", this.onEndedCallback.bind(this));
    }
  }

  playAudio(track: AudioTrack) {
    this.stopAudio();
    this.currentPlayingTrack = track;
    this.audios[track].play();
    this.audios[track].addEventListener("ended", this.onEndedCallback.bind(this));
  }
}

const audio = new SettlersAudio();
export { audio };
