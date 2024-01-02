# Yoga with Adriene - daily yoga video from her YouTube playlist.
# Script opens Chrome browser, opens monthly playlist and picks up video with today's date.
# It turns off captions with "c" button and goes to full screen with "f" key.
# Replace JSON and it should work!

import datetime
import webbrowser
from google.oauth2.service_account import Credentials
from googleapiclient.discovery import build
import pyautogui
import os
import time
import subprocess

# Replace the channel ID and path to the JSON file with your own
channel_id = "UCFKE7WVJfvaHW5q283SxchA"
creds = Credentials.from_service_account_file(r"alpine-realm-381711-182fcef9362c.json")

# Authenticate with the YouTube Data API
youtube = build('youtube', 'v3', credentials=creds)

# Get the playlists for the selected channel
playlists = []
next_page_token = None
while True:
    request = youtube.playlists().list(
        part='snippet',
        channelId=channel_id,
        maxResults=50,
        pageToken=next_page_token
    )
    response = request.execute()
    playlists += response['items']
    next_page_token = response.get('nextPageToken')
    if not next_page_token:
        break

# Filter the playlists based on the current month
now = datetime.datetime.now()
month = now.strftime("%B")
filtered_playlists = []
for playlist in playlists:
    if month.lower() in playlist['snippet']['title'].lower():
        filtered_playlists.append(playlist)

# Select the first playlist from the filtered playlists
selected_playlist = filtered_playlists[0]['id']

# Get the playlist items for the selected playlist
playlist_items = []
next_page_token = None
while True:
    request = youtube.playlistItems().list(
        part='snippet',
        playlistId=selected_playlist,
        maxResults=50,
        pageToken=next_page_token
    )
    response = request.execute()
    playlist_items += response['items']
    next_page_token = response.get('nextPageToken')
    if not next_page_token:
        break

# Calculate the index of the video for the current day
now = datetime.datetime.now()
video_index = (now.day - 1) % len(playlist_items)

# Get the video ID for the selected video
video_id = playlist_items[video_index]['snippet']['resourceId']['videoId']

# Build the URL for the selected video
video_url = f"https://www.youtube.com/watch?v={video_id}"

# Open the selected video in a new browser window
webbrowser.open_new(video_url)

# Command to move the browser window to the right side of the primary display
moveWindow = 'xdotool search --onlyvisible --class "browser" windowmove --relative -- 1920 0'
setAudio = 'pactl set-default-sink alsa_output.pci-0000_01_00.1.hdmi-stereo-extra1'
# Wait for the page to load
time.sleep(5)

# Move browser to the right and set audio device to HDMI2
combined_command = f"{moveWindow}; {setAudio}"
#subprocess.run(combined_command, shell=True)

# Press "f" to enter full screen
pyautogui.press('f')

# Press "c" to disable captions
pyautogui.press('c')
