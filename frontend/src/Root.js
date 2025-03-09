import React, { useEffect, useState } from 'react';
import { Container, Typography, Box } from '@mui/material';

const Root = ({ url }) => {
  const [data, setData] = useState(null);

  useEffect(() => {
    const params = new URLSearchParams(window.location.search);
    const user = params.get('user');

    fetch(`${url}/?user=${user}`)
      .then(response => response.text())
      .then(data => setData(data))
      .catch(error => console.error('Error fetching data:', error));
  }, [url]);

  return (
    <Container maxWidth="sm">
      <Box sx={{ mt: 4 }}>
        <Typography variant="h2" component="h1" gutterBottom>
          {data}
        </Typography>
      </Box>
    </Container>
  );
};

export default Root;