describe('Registeration to login navigation', () => {
  it('should nigivate to /login when click Login button', () => {
    // visit registration page
    cy.visit('/registration');

    // click login button
    cy.get('[data-test="login-button"]').click(); 

    // check if the url change to /login
    cy.url().should('include', '/login'); 
  })
})