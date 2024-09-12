import React, { useState } from 'react';

const url = process.env.REACT_APP_BACKEND_URL;

const VideoAdd = () => {
  const [id, setId] = useState('');
  const [title, setTitle] = useState('');

  const handleSubmit = (e) => {
    e.preventDefault();
    fetch(`${url}/video?id=${id}&title=${title}`, {
      method: 'POST',
    })
      .then(response => {
        if (response.ok) {
          alert('Video added successfully');
        } else {
          alert('Failed to add video');
        }
      })
      .catch(error => console.error('Error adding video:', error));
  };

  return (
    <div>
      <h1>Add Video</h1>
      <form onSubmit={handleSubmit}>
        <div>
          <label>ID:</label>
          <input type="text" value={id} onChange={(e) => setId(e.target.value)} />
        </div>
        <div>
          <label>Title:</label>
          <input type="text" value={title} onChange={(e) => setTitle(e.target.value)} />
        </div>
        <button type="submit">Add Video</button>
      </form>
    </div>
  );
};

export default VideoAdd;
