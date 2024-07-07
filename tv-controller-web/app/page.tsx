'use client'

import Image from "next/image";
import ConfigPanel, {Config} from "@/app/config/config";
import {useEffect, useState} from "react";
import {getConfig, GetConfigResponse} from "@/app/vlccontrol";

export default function Home() {
  const [ cfg, setCfg ] = useState<GetConfigResponse>({
      playlists_available: ['test1.m3u', 'test2.m3u'],
      selected_playlist: '',
      shuffle: true,
      loop: true,
      vlc_path: ''
  });

  return (
      <main>
        <h1>TV Control</h1>
        <ConfigPanel cfg={cfg} />
      </main>
  );
}
