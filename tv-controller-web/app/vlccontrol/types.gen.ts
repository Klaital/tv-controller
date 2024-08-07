// This file is auto-generated by @hey-api/openapi-ts

export type CurrentPlaybackSettings = {
    playlists_available?: Array<(string)>;
    selected_playlist?: string;
    shuffle?: boolean;
    loop?: boolean;
    vlc_path?: string;
};

export type NewPlaybackSettings = {
    selected_playlist?: string;
    shuffle?: boolean;
    loop?: boolean;
    vlc_path?: string;
};

export type SelectPlaylistRequest = {
    playlist?: string;
};

export type PausePlaybackResponse = unknown;

export type PausePlaybackError = unknown;

export type TrackAheadResponse = unknown;

export type TrackAheadError = unknown;

export type TrackBackResponse = unknown;

export type TrackBackError = unknown;

export type GetConfigResponse = CurrentPlaybackSettings;

export type GetConfigError = unknown;

export type SetConfigData = {
    body?: NewPlaybackSettings;
};

export type SetConfigResponse = CurrentPlaybackSettings;

export type SetConfigError = unknown;

export type ToggleShuffleResponse = unknown;

export type ToggleShuffleError = unknown;

export type ToggleLoopResponse = unknown;

export type ToggleLoopError = unknown;

export type SelectPlaylistData = {
    /**
     * Specify a playlist to use
     */
    body?: SelectPlaylistRequest;
};

export type SelectPlaylistResponse = unknown;

export type SelectPlaylistError = unknown;

export type $OpenApiTs = {
    '/ctrl/pause': {
        put: {
            res: {
                /**
                 * Success
                 */
                '200': unknown;
            };
        };
    };
    '/ctrl/trackahead': {
        put: {
            res: {
                /**
                 * Success
                 */
                '200': unknown;
            };
        };
    };
    '/ctrl/trackback': {
        put: {
            res: {
                /**
                 * Success
                 */
                '200': unknown;
            };
        };
    };
    '/cfg': {
        get: {
            res: {
                /**
                 * success
                 */
                '200': CurrentPlaybackSettings;
            };
        };
        put: {
            req: SetConfigData;
            res: {
                /**
                 * success
                 */
                '200': CurrentPlaybackSettings;
            };
        };
    };
    '/cfg/shuffle': {
        put: {
            res: {
                /**
                 * success
                 */
                '200': unknown;
            };
        };
    };
    '/cfg/loop': {
        put: {
            res: {
                /**
                 * success
                 */
                '200': unknown;
            };
        };
    };
    '/cfg/playlist': {
        put: {
            req: SelectPlaylistData;
            res: {
                /**
                 * success
                 */
                '200': unknown;
            };
        };
    };
};