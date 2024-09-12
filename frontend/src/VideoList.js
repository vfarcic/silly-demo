import React, { useEffect, useState } from 'react';

const url = process.env.REACT_APP_BACKEND_URL;

const VideoList = () => {
  const [videos, setVideos] = useState([]);

  useEffect(() => {
    fetch(`${url}/videos`)
      .then(response => response.json())
      .then(data => setVideos(data))
      .catch(error => console.error('Error fetching videos:', error));
  }, []);

  return (
    <div>
      <h1>List Videos</h1>
      <ul>
        {videos.map(video => (
          <li key={video.id}>
            {video.title} (ID: {video.id})
          </li>
        ))}
      </ul>
    </div>
  );
};

export default VideoList;
