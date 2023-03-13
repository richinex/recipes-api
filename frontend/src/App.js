import React, { useEffect, useState } from 'react';
import './App.css';
import Recipes from './Recipes';
import Navbar from './Navbar';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';


function App(props) {
  const [recipes, setRecipes] = useState([]);

  useEffect(() => {
    getRecipes();
  }, []);

  function getRecipes() {
    fetch('http://localhost:8080/recipes')
      .then(response => response.json())
      .then(data => setRecipes(data));
  }

  return (
    <div>
      <Router>
        <Navbar />
        <Routes>
        <Route path="/" element={<Recipes recipes={recipes} />} />
        </Routes>
      </Router>
    </div>
  );
}

export default App;