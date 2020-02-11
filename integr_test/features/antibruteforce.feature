Feature: antibruteforce
  In order to use antibruteforce
  As an GRPC client
  I need to be able to send grpc requests

  Scenario Outline: should check for bruteforce
    Given login is "<login>"
    And password is "<password>"
    And ip is "<ip>"
    When I call grpc method "Check"
    Then response error should be "<error>"

    Examples:
      | login  | password | ip      | error              |
      | test1 | pass1    | 1.2.3.4 | nil                |
      | test1 | pass1    | 1.2.3.4 | nil                |
      | test1 | pass1    | 1.2.3.4 | bucket is overflow |
      | test1 | pass2    | 1.2.3.4 | bucket is overflow |

  Scenario Outline: need to reset bucket
    Given login is "<login>"
    And ip is "<ip>"
    And password is "<password>"

    When I call grpc method "<method>"
    Then response error should be "<error>"

    Examples:
      | method | login  | password | ip      | error              |
      | Auth  | test1 | pass2    | 1.2.3.4 | bucket is overflow |
      | Drop  | test1 |          | 1.2.3.4 | nil                |
      | Auth  | test1 | pass2    | 1.2.3.4 | nil                |

  Scenario Outline: need to add subnet to blacklist
    Given login is "<login>"
    And ip is "<ip>"
    And password is "<password>"
    And subnet is "192.168.0.0/25"

    When I call grpc method "<method>"
    Then response error should be "<error>"

    Examples:
      | method         | login  | password | ip           | error            |
      | Auth           | test3 | pass3    | 172.16.1.10  | nil              |
      | AddToBlacklist |        |          |              | nil              |
      | Auth           | test3 | pass3    | 172.16.1.10  | ip in black list |

  Scenario Outline: need to remove subnet from blacklist
    Given login is "<login>"
    And ip is "<ip>"
    And password is "<password>"
    And subnet is "192.168.0.0/25"

    When I call grpc method "<method>"
    Then response error should be "<error>"

    Examples:
      | method          | login  | password | ip           | error            |
      | Check           | test4 | pass4    | 192.168.0.30 | ip in black list |
      | BlacklistRemove |        |          |              | nil              |
      | Check           | test4 | pass4    | 192.168.0.30 | nil              |

  Scenario Outline: need to add subnet to whitelist
    Given login is "<login>"
    And ip is "<ip>"
    And password is "<password>"
    And subnet is "192.168.0.0/25"

    When I call grpc method "<method>"
    Then response error should be "<error>"

    Examples:
      | method       | login  | password | ip           | error              |
      | Check        | test5 | pass5    | 192.168.0.30 | nil                |
      | Check        | test5 | pass5    | 192.168.0.30 | nil                |
      | Check        | test5 | pass5    | 192.168.0.30 | bucket is overflow |
      | WhitelistAdd |        |          |              | nil                |
      | Check        | test5 | pass5    | 192.168.0.30 | nil                |

  Scenario Outline: need to remove subnet from whitelist
    Given login is "<login>"
    And ip is "<ip>"
    And password is "<password>"
    And subnet is "192.168.0.0/25"

    When I call grpc method "<method>"
    Then response error should be "<error>"

    Examples:
      | method         | login | password | ip           | error              |
      | Auth           | test6 | pass6    | 172.16.1.10 | nil                |
      | Auth           | test6 | pass6    | 172.16.1.10 | nil                |
      | Auth           | test6 | pass6    | 172.16.1.10 | nil                |
      | RemoveFromWhitelist    |          |              |              | nil                |
      | Auth           | test6 | pass6    | 172.16.1.10 | nil                |
      | Auth           | test6 | pass6    | 172.16.1.10 | nil                |
      | Auth           | test6 | pass6    | 172.16.1.10 | bucket is overflow |