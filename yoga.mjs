const { google } = require('googleapis');
import open from 'open';
const fs = require('fs');
const path = require('path');

const secretFile = path.join(__dirname, 'alpine-realm-381711-182fcef9362c.json');
const credentials = JSON.parse(fs.readFileSync(secretFile, 'utf8'));

const jwtClient = new google.auth.JWT(
  credentials.client_email,
  null,
  credentials.private_key,
  ['https://www.googleapis.com/auth/youtube.readonly'],
  null
);

google.options({ auth: jwtClient });

const youtube = google.youtube('v3');

async function openVideoInBrowser() {
  try {
    const playlistsResponse = await youtube.playlists.list({
      part: 'snippet',
      channelId: 'UCFKE7WVJfvaHW5q283SxchA',
      maxResults: 50,
    });

    const currentMonth = new Date().toLocaleString('default', { month: 'long' });
    let selectedPlaylist = playlistsResponse.data.items.find(playlist =>
      playlist.snippet.title.includes(currentMonth)
    );

    if (!selectedPlaylist) {
      console.log('No playlist found for the current month');
      return;
    }

    const playlistItemsResponse = await youtube.playlistItems.list({
      part: 'snippet',
      playlistId: selectedPlaylist.id,
      maxResults: 50,
    });

    const videoIndex = (new Date().getDate() - 1) % playlistItemsResponse.data.items.length;
    const videoId = playlistItemsResponse.data.items[videoIndex].snippet.resourceId.videoId;
    const videoURL = `https://www.youtube.com/watch?v=${videoId}`;

    await open(videoURL);
  } catch (err) {
    console.error('Error: ', err);
  }
}

openVideoInBrowser();
