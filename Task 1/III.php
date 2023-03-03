interface OrderInterface {
  public function createOrder($orderData);
  public function getOrder($orderId);
  public function updateOrder($orderId, $orderData);
  public function addItem($productId, $quantity);
  public function removeItem($productId);
  public function placeOrder($userId, $orderId, $shippingAddress, $billingAddress, $paymentMethod);
  public function deleteOrder($orderId);
}

interface CardInterface {
  public function addToCard($CardData);
  public function getCard($CardId);
  public function updateCard($CardId, $CardData);
  public function deleteCard($CardId);
}

interface ProductInterface {
  public function getProducts($filterData);
  public function getProduct($productId);
}

class Order implements OrderInterface {
  public function createOrder($orderData) {}
  public function getOrder($orderId) {}
  public function addItem($productId, $quantity) {}
  public function removeItem($productId) {}
  public function updateOrder($orderId, $orderData) {}
  public function placeOrder($userId, $orderId, $shippingAddress, $billingAddress, $paymentMethod) {}
  public function deleteOrder($orderId) {}
}

class Card implements CardInterface {
  public function addToCard($CardData) {}
  public function getCard($CardId) {}
  public function updateCard($CardId, $CardData) {}
  public function deleteCard($CardId) {}
}

class Product implements ProductInterface {
  public function getProducts($filterData) {}
  public function getProduct($productId) {}
}