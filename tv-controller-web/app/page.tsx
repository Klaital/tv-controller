'use client'

import Image from "next/image";
import ConfigPanel, {Config} from "@/app/config/config";
import React, {useEffect, useState} from "react";
import {getConfig, GetConfigResponse, selectPlaylist} from "@/app/vlccontrol";
import {createClient} from "@hey-api/client-fetch";


createClient({
    baseUrl:'http://localhost:8080',
})

export default function Home() {
  const [ cfg, setCfg ] = useState<GetConfigResponse>({
      playlists_available: ['test1.m3u', 'test2.m3u'],
      selected_playlist: '',
      shuffle: true,
      loop: true,
      vlc_path: ''
  });

  function chooseNewPlaylist(newPlaylist: string) {
      // e.preventDefault();
      console.log("Selected Playlist: " + newPlaylist);

      // TODO: send request to backend to change the active playlist
      selectPlaylist({
          body: {
              playlist: newPlaylist,
          },
      })
          .then((resp) => {
              console.log("active playlist updated");
              setCfg({
                  playlists_available: cfg.playlists_available,
                  selected_playlist: newPlaylist,
                  shuffle: cfg.shuffle,
                  loop: cfg.loop,
                  vlc_path: cfg.vlc_path,
              })
          })
          .catch(e => {
              console.log(e);
          })
  }

  function loadCfg() {
      getConfig()
          .then((resp) => {
              setCfg({
                  playlists_available: resp.data?.playlists_available,
                  selected_playlist: resp.data?.selected_playlist,
                  shuffle: resp.data?.shuffle,
                  loop: resp.data?.loop,
                  vlc_path: resp.data?.vlc_path,
              })
          })
          .catch(e => {
              console.log(e)
          })
  }

    useEffect(() => {
        loadCfg();
    }, []);
  return (
      <main>
        <h1>TV Control</h1>
        <ConfigPanel cfg={cfg} selectPlaylist={chooseNewPlaylist} />
      </main>
  );
}
