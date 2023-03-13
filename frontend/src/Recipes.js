import React from 'react';

function Recipes({ recipes }) {

  function Recipe(props) {
    return (
      <div className="recipe">
        <h4>{props.recipe.name}</h4>
        <ul>
          {props.recipe.ingredients.map((ingredient, index) => (
            <li key={index}>{ingredient}</li>
          ))}
        </ul>
      </div>
    );
  }

  return (
    <div>
      {recipes.map((recipe, index) => (
        <Recipe recipe={recipe} key={index} />
      ))}
    </div>
  );
}

export default Recipes;
