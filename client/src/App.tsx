import './App.css';

import { instance } from './axios';

function App() {
  const handleClick = () => {
    instance
      .get('/')
      .then(res => {
        console.log(res);
      })
      .catch(err => {
        console.log(err);
      });
  };

  return (
    <button type="button" onClick={handleClick}>
      Click
    </button>
  );
}

export default App;
