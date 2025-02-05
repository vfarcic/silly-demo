import React, { useState } from 'react';
import { TextField, Button, Container, Typography, Box } from '@mui/material';

const url = process.env.REACT_APP_BACKEND_URL;

const VideoAdd = () => {
  const [id, setId] = useState('');
  const [title, setTitle] = useState('');

  const handleSubmit = (event) => {
    event.preventDefault();
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
    <Container maxWidth="sm">
      <Box sx={{ mt: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          Add Video
        </Typography>
        <form onSubmit={handleSubmit}>
          <TextField
            label="ID"
            variant="outlined"
            fullWidth
            margin="normal"
            value={id}
            onChange={(e) => setId(e.target.value)}
          />
          <TextField
            label="Title"
            variant="outlined"
            fullWidth
            margin="normal"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
          />
          <Button type="submit" variant="contained" color="primary" fullWidth>
            Add Video
          </Button>
        </form>
      </Box>
    </Container>
  );
};

export default VideoAdd;
