import './App.css';
import React from 'react';
import VideoList from './VideoList';
import AddVideo from './AddVideo';

function App() {
  return (
    <div className="App">
      <AddVideo />
      <VideoList />
    </div>
  );
}

export default App;
