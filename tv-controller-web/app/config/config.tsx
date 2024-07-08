import {GetConfigResponse} from "@/app/vlccontrol";
import React from "react";

function PlaylistButton(props: {playlistId: string, selectPlaylist: (newPlaylist: string) => void}) {
    return (
        <a
            key={props.playlistId}
            onClick={(e: React.MouseEvent<HTMLElement>) => {
                // e.preventDefault();
                props.selectPlaylist(props.playlistId);
            }}
            className="playlist inline-block bg-gray-200 rounded-full px-3 py-1 text-sm font-semibold text-gray-700 mr-2 mb-2"
        >
            {props.playlistId}
        </a>
    )
}

export default function ConfigPanel(props: {
    cfg: GetConfigResponse
    selectPlaylist: (newPlaylist: string) => void
}) {

    return <>
        <div className="playlist-selection">
            {props.cfg.playlists_available?.map((playlist_name: string) => (
                <PlaylistButton key={playlist_name} playlistId={playlist_name} selectPlaylist={props.selectPlaylist} />
            ))}
        </div>
    </>
}
